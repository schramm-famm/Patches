package websockets

import (
	"log"
	"patches/protocol"
	"time"

	"github.com/gorilla/websocket"
	gorillaws "github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Client manages a WebSocket connection with a client.
type Client struct {
	userID         int64
	conversationID int64
	caret          protocol.Caret
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
		caret:          protocol.Caret{Start: 0, End: 0},
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

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(
					gorillaws.CloseMessage,
					gorillaws.FormatCloseMessage(gorillaws.CloseGoingAway, "Going away"),
				)
				return
			}

			err := c.conn.WriteMessage(gorillaws.TextMessage, message)
			if err != nil {
				log.Print("Failed to write message to WebSocket: ", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Print("Failed to write ping message to WebSocket: ", err)
				return
			}
		}
	}
}
