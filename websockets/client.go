package websockets

import (
	"log"

	gorillaws "github.com/gorilla/websocket"
)

// Client manages a WebSocket connection with a client.
type Client struct {
	userID         int64
	conversationID int64
	conn           *gorillaws.Conn
	broker         *Broker
	broadcast      chan<- *Message
	send           <-chan []byte
}

// NewClient creates a new Client struct.
func NewClient(userID int64, conn *gorillaws.Conn, broadcast chan<- *Message, broker *Broker) *Client {
	return &Client{
		userID:    userID,
		conn:      conn,
		broadcast: broadcast,
		broker:    broker,
		send:      make(chan []byte),
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

		c.broadcast <- &Message{data, c}
	}
}

// write sends messages to the WebSocket connection whenever new messages are
// sent into the Client's channel.
func (c *Client) write() {
	for message := range c.send {
		err := c.conn.WriteMessage(gorillaws.TextMessage, message)
		if err != nil {
			return
		}
	}
	c.conn.WriteMessage(gorillaws.CloseMessage, []byte{})
}
