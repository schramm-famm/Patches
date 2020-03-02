package websockets

import (
	"log"

	gorillaws "github.com/gorilla/websocket"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Conversation manages all WebSocket connections in a single conversation.
type Conversation struct {
	conversationID int64
	doc            string
	clients        map[*Client]bool

	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

// Message stores the content and sender of a message.
type Message struct {
	content []byte
	sender  *Client
}

var dmp *diffmatchpatch.DiffMatchPatch = diffmatchpatch.New()

// NewConversation creates a new Conversation struct.
func NewConversation(conversationID int64, doc string) *Conversation {
	return &Conversation{
		conversationID: conversationID,
		doc:            doc,
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan *Message),
	}
}

func (c *Conversation) deleteClient(client *Client) {
	delete(c.clients, client)
	close(client.send)
	log.Printf("Deregistered a client in conversation %d (%d active)", c.conversationID, len(c.clients))
}

// Run waits on a Conversation's three channels for clients to be added, clients
// to be removed, and messages to be broadcast. Only one of these operations may
// be performed at a time.
func (c *Conversation) Run() {
	for {
		select {
		case client := <-c.register:
			err := client.conn.WriteMessage(gorillaws.TextMessage, []byte(c.doc))
			if err != nil {
				log.Print("Failed to send initial conversation data: ", err)
				close(client.send)
				continue
			}
			c.clients[client] = true
			log.Printf("Registered a client in conversation %d (%d active)", c.conversationID, len(c.clients))

		case client := <-c.unregister:
			if _, ok := c.clients[client]; ok {
				c.deleteClient(client)
			} else {
				log.Printf("Attempted to deregister an inactive client in conversation %d", c.conversationID)
			}

		case message, ok := <-c.broadcast:
			if !ok {
				close(c.register)
				close(c.unregister)
				log.Printf("Shutting down conversation %d", c.conversationID)
				return
			}

			if _, ok := c.clients[message.sender]; !ok {
				continue
			}

			patches, err := dmp.PatchFromText(string(message.content))
			if err != nil {
				log.Printf("Failed to create patch from %s: %v", message.content, err)
				c.deleteClient(message.sender)
				continue
			}
			newDoc, okList := dmp.PatchApply(patches, c.doc)
			for i, ok := range okList {
				if !ok {
					log.Printf("Failed to apply patch %+v", patches[i])
					c.deleteClient(message.sender)
					continue
				}
			}
			c.doc = newDoc

			for client := range c.clients {
				if client != message.sender {
					client.send <- message.content
				}
			}
		}
	}
}
