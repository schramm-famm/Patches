package websockets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"patches/models"
	"strconv"
	"sync"
	"time"

	gorillaws "github.com/gorilla/websocket"
)

const (
	authRoute    = "/heimdall/v1/token/auth"
	memberRoute  = "/ether/v1/conversations/%d/users/%d"
	contentRoute = "/ether/v1/conversations/%d/content"
)

var (
	heimdallHost = os.Getenv("HEIMDALL_HOST")
	etherHost = os.Getenv("ETHER_HOST")
)

// ConvoData represents a conversation and its associated clients.
type ConvoData struct {
	conversation *Conversation
	clients      map[*Client]bool
}

// Broker is the single entrypoint for registering and unregistering all
// client connections.
type Broker struct {
	sync.Mutex
	active     map[int64]*ConvoData
	db         models.Datastore
	httpClient *http.Client
}

// NewBroker creates a new Broker struct.
func NewBroker(db models.Datastore, httpClient *http.Client) *Broker {
	return &Broker{
		active:     make(map[int64]*ConvoData),
		db:         db,
		httpClient: httpClient,
	}
}

// register adds a client connection to an active conversation.
func (b *Broker) register(member *models.UserConversationMapping, conn *gorillaws.Conn) (*Client, error) {
	b.Lock()
	defer b.Unlock()

	cd, ok := b.active[member.ConversationID]
	if !ok {
		// If this is the first client connection in this conversation, then get
		// the HTML content and create a new Conversation struct to manage the
		// conversation
		content, err := b.getConversationContent(member.UserID, member.ConversationID)
		if err != nil {
			return nil, err
		}

		cd = &ConvoData{
			conversation: NewConversation(member.ConversationID, content, b),
			clients:      make(map[*Client]bool),
		}
		go cd.conversation.Run()
		b.active[member.ConversationID] = cd
	}

	client := NewClient(member.UserID, conn, cd.conversation, b)
	cd.clients[client] = true
	cd.conversation.register <- client
	return client, nil
}

// unregister removes a client connection from its associated conversation.
func (b *Broker) unregister(client *Client) {
	b.Lock()
	defer b.Unlock()

	conversationID := client.conversation.conversationID
	cd, ok := b.active[conversationID]
	if ok {
		client.conversation.unregister <- client
		delete(cd.clients, client)
		if len(cd.clients) == 0 {
			delete(b.active, conversationID)
			close(client.conversation.broadcast)
		}
	}
}

// validate token checks with Heimdall whether a token is authentic and returns
// the embedded user ID if it is.
func (b *Broker) validateToken(token string) (int64, error) {
	reqBody, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		return -1, err
	}

	res, err := b.httpClient.Post(
		"http://"+heimdallHost+authRoute,
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return -1, err
	}

	if res.StatusCode != http.StatusNotFound && res.StatusCode != http.StatusOK {
		return -1, errors.New("Failed to validate token")
	} else if res.StatusCode == http.StatusNotFound {
		return -1, errors.New("Token is invalid")
	}

	resBody := struct {
		UserID int64 `json:"user_id"`
	}{}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return -1, err
	}

	return resBody.UserID, nil
}

func (b *Broker) getConversationMember(userID, conversationID int64) (*models.UserConversationMapping, error) {
	req, err := http.NewRequest("GET", "http://"+etherHost+fmt.Sprintf(memberRoute, conversationID, userID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-ID", strconv.FormatInt(userID, 10))
	res, err := b.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusNotFound && res.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to get conversation member")
	} else if res.StatusCode == http.StatusNotFound {
		return nil, errors.New("Conversation/member not found")
	}

	resBody := models.UserConversationMapping{}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, err
	}

	return &resBody, nil
}

// getConversationContent gets the HTML content of a conversation from Ether.
func (b *Broker) getConversationContent(userID, conversationID int64) (string, error) {
	req, err := http.NewRequest("GET", "http://"+etherHost+fmt.Sprintf(contentRoute, conversationID), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-ID", strconv.FormatInt(userID, 10))
	res, err := b.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusNotFound && res.StatusCode != http.StatusOK {
		return "", errors.New("Failed to get conversation content")
	} else if res.StatusCode == http.StatusNotFound {
		return "", errors.New("Conversation not found")
	}

	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// StartClient authenticates a WebSocket connection before registering the
// connection and starting goroutines for reading to and writing from the
// connection.
func (b *Broker) StartClient(conversationID int64, conn *gorillaws.Conn) {
	// Wait for client to send token through the WebSocket connection
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		netErr, ok := err.(net.Error)
		if ok && netErr.Timeout() {
			log.Print("Timed out waiting for a token from the client")
		}
		if gorillaws.IsUnexpectedCloseError(err, gorillaws.CloseGoingAway, gorillaws.CloseAbnormalClosure) {
			log.Printf("WebSocket closed unexpectedly: %v", err)
		}
		conn.Close()
		return
	}
	conn.SetReadDeadline(time.Time{})

	// Verify with Heimdall that the token is authentic
	token := string(message)
	userID, err := b.validateToken(token)
	if err != nil {
		log.Print("Failed to validate token: ", err)
		conn.Close()
		return
	}

	member, err := b.getConversationMember(userID, conversationID)
	if err != nil {
		log.Printf(
			"Failed to validate conversation member (user: %d, conversation: %d): %v",
			userID,
			conversationID,
			err,
		)
		conn.Close()
		return
	}

	// Create client struct with the user ID and start reading/writing patches
	client, err := b.register(member, conn)
	if err != nil {
		log.Printf("Failed to create a new client (user: %d, conversation: %d): %v", userID, conversationID, err)
		conn.Close()
		return
	}
	go client.write()
	go client.read()
}
