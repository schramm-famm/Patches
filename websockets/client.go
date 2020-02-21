package websockets

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	gorillaws "github.com/gorilla/websocket"
)

const (
	authRoute = "/heimdall/v1/token/auth"
)

var heimdallHost = os.Getenv("HEIMDALL_HOST")

type Client struct {
	userID       int64
	conn         *gorillaws.Conn
	conversation *Conversation
	send         chan []byte
}

func NewClient(userID int64, conn *gorillaws.Conn, conversation *Conversation) *Client {
	return &Client{
		userID:       userID,
		conn:         conn,
		conversation: conversation,
		send:         make(chan []byte),
	}
}

func (c *Client) validateToken(token string) (int64, error) {
	reqBody, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		return -1, err
	}

	res, err := c.conversation.broker.httpClient.Post(
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

func (c *Client) read() {
	defer func() {
		c.conversation.unregister <- c
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if gorillaws.IsUnexpectedCloseError(err, gorillaws.CloseGoingAway, gorillaws.CloseAbnormalClosure) {
				log.Printf("WebSocket closed unexpectedly: %v", err)
			}
			break
		}

		c.conversation.broadcast <- &Message{message, c}
	}
}

func (c *Client) write() {
	for {
		message, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(gorillaws.CloseMessage, []byte{})
			return
		}

		err := c.conn.WriteMessage(gorillaws.TextMessage, message)
		if err != nil {
			return
		}
	}
}

func (c *Client) Run() {
	defer func() {
		c.conn.Close()
		c.conversation.unregister <- c
	}()

	// Wait for client to send token through the WebSocket connection
	c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, message, err := c.conn.ReadMessage()
	if err != nil {
		netErr, ok := err.(net.Error)
		if ok && netErr.Timeout() {
			log.Print("Timed out waiting for a token from the client")
		}
		if gorillaws.IsUnexpectedCloseError(err, gorillaws.CloseGoingAway, gorillaws.CloseAbnormalClosure) {
			log.Printf("WebSocket closed unexpectedly: %v", err)
		}
		return
	}
	c.conn.SetReadDeadline(time.Time{})

	// Verify with Heimdall that the token is authentic
	token := string(message)
	userID, err := c.validateToken(token)
	if err != nil {
		log.Print("Failed to validate token: ", err)
		return
	}

	// Update the client struct with the user ID and start reading/writing patches
	c.userID = userID
	c.conversation.activate <- c
	go c.write()
	c.read()
}
