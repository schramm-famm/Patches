package websockets

import (
	"log"

	gorillaws "github.com/gorilla/websocket"
)

// Client manages a WebSocket connection with a client.
type Client struct {
	userID       int64
	conn         *gorillaws.Conn
	conversation *Conversation
	broker       *Broker
	send         chan []byte
}

// NewClient creates a new Client struct.
func NewClient(userID int64, conn *gorillaws.Conn, conversation *Conversation, broker *Broker) *Client {
	return &Client{
		userID:       userID,
		conn:         conn,
		conversation: conversation,
		broker:       broker,
		send:         make(chan []byte),
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

		c.conversation.broadcast <- &Message{data, c}
	}
}

// write sends messages to the WebSocket connection whenever new messages are
// sent into the Client's channel.
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
