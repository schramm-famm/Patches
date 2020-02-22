package websockets

import (
	"log"

	gorillaws "github.com/gorilla/websocket"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Conversation struct {
	clients        map[*Client]bool
	broker         *Broker
	conversationID int64
	doc            *string

	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

type Message struct {
	content []byte
	sender  *Client
}

var dmp *diffmatchpatch.DiffMatchPatch = diffmatchpatch.New()

func NewConversation(conversationID int64, broker *Broker) *Conversation {
	doc := ""
	return &Conversation{
		clients:        make(map[*Client]bool),
		broker:         broker,
		conversationID: conversationID,
		doc:            &doc,
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan *Message),
	}
}

func (c *Conversation) Run() {
	for {
		select {
		case client := <-c.register:
			c.clients[client] = true
			err := client.conn.WriteMessage(gorillaws.TextMessage, []byte(*c.doc))
			if err != nil {
				return
			}
			log.Printf("Registered a client in conversation %d (%d active)", c.conversationID, len(c.clients))

		case client := <-c.unregister:
			if _, ok := c.clients[client]; ok {
				delete(c.clients, client)
				close(client.send)
				log.Printf("Deregistered a client in conversation %d (%d active)", c.conversationID, len(c.clients))
			}

		case message, ok := <-c.broadcast:
			if !ok {
				close(c.register)
				close(c.unregister)
				log.Printf("Shutting down conversation %d", c.conversationID)
				return
			}

			patches, err := dmp.PatchFromText(string(message.content))
			if err != nil {
				return
			}
			*c.doc, _ = dmp.PatchApply(patches, *c.doc)

			for client := range c.clients {
				if client != message.sender {
					client.send <- message.content
				}
			}
		}
	}
}
