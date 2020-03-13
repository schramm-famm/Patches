package websockets

import (
	"log"
	"patches/models"

	gorillaws "github.com/gorilla/websocket"
)

// Client manages a WebSocket connection with a client.
type Client struct {
	userID         int64
	conversationID int64
	caret          models.Caret
	conn           *gorillaws.Conn
	broker         *Broker
	broadcast      chan<- *BroadcastMessage
	send           chan []byte
}

// NewClient creates a new Client struct.
func NewClient(
	userID int64,
	conversationID int64,
	conn *gorillaws.Conn,
	broker *Broker,
	broadcast chan<- *BroadcastMessage,
) *Client {
	return &Client{
		userID:         userID,
		conversationID: conversationID,
		caret:          models.Caret{Start: 0, End: 0},
		conn:           conn,
		broadcast:      broadcast,
		broker:         broker,
		send:           make(chan []byte),
	}
}

// read consumes messages from the WebSocket connection and sends them to the
// associated conversation to be broadcast.
func (c *Client) read() {
	defer func() {
		c.conn.Close()
		c.broker.unregister(c)
	}()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if gorillaws.IsUnexpectedCloseError(err, gorillaws.CloseGoingAway, gorillaws.CloseAbnormalClosure) {
				log.Printf("WebSocket closed unexpectedly: %v", err)
			}
			break
		}

		c.broadcast <- &BroadcastMessage{data, c}
	}
}

// write sends messages to the WebSocket connection whenever new messages are
// sent into the Client's channel.
func (c *Client) write() {
	defer func() {
		c.conn.WriteMessage(gorillaws.CloseMessage, []byte{})
	}()

	for message := range c.send {
		err := c.conn.WriteMessage(gorillaws.TextMessage, message)
		if err != nil {
			log.Print("Failed to write message to WebSocket: ", err)
			return
		}
	}
}
