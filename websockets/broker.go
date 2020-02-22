package websockets

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"patches/models"
	"sync"
	"time"

	gorillaws "github.com/gorilla/websocket"
)

const (
	authRoute = "/heimdall/v1/token/auth"
)

var heimdallHost = os.Getenv("HEIMDALL_HOST")

type ConvoData struct {
	conversation *Conversation
	clients      map[*Client]bool
}

type Broker struct {
	sync.Mutex
	active     map[int64]*ConvoData
	db         models.Datastore
	httpClient *http.Client
}

func NewBroker(db models.Datastore, httpClient *http.Client) *Broker {
	return &Broker{
		active:     make(map[int64]*ConvoData),
		db:         db,
		httpClient: httpClient,
	}
}

func (b *Broker) register(userID, conversationID int64, conn *gorillaws.Conn) *Client {
	b.Lock()
	defer b.Unlock()

	cd, ok := b.active[conversationID]
	if !ok {
		cd = &ConvoData{
			conversation: NewConversation(conversationID, b),
			clients:      make(map[*Client]bool),
		}
		go cd.conversation.Run()
		b.active[conversationID] = cd
	}

	client := NewClient(userID, conn, cd.conversation, b)
	cd.clients[client] = true
	cd.conversation.register <- client
	return client
}

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

	// Create client struct with the user ID and start reading/writing patches
	client := b.register(userID, conversationID, conn)
	go client.write()
	go client.read()
}
