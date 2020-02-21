package websockets

import (
	"log"

	gorillaws "github.com/gorilla/websocket"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Conversation struct {
	pendingClients map[*Client]bool
	activeClients  map[*Client]bool
	broker         *Broker
	conversationID int64
	doc            *string

	register   chan *Client
	activate   chan *Client
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
		pendingClients: make(map[*Client]bool),
		activeClients:  make(map[*Client]bool),
		broker:         broker,
		conversationID: conversationID,
		doc:            &doc,
		register:       make(chan *Client),
		activate:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan *Message),
	}
}

func (c *Conversation) Run() {
	for {
		select {
		case client := <-c.register:
			c.pendingClients[client] = true
			log.Printf("Registered a client (%d pending, %d active)", len(c.pendingClients), len(c.activeClients))

		case client := <-c.activate:
			if _, ok := c.pendingClients[client]; ok {
				delete(c.pendingClients, client)
				c.activeClients[client] = true
				err := client.conn.WriteMessage(gorillaws.TextMessage, []byte(*c.doc))
				if err != nil {
					return
				}
				log.Printf("Activated a client (%d pending, %d active)", len(c.pendingClients), len(c.activeClients))
			}

		case client := <-c.unregister:
			if _, ok := c.pendingClients[client]; ok {
				delete(c.pendingClients, client)
				close(client.send)
				log.Printf("Deregistered a client (%d clients, %d active)", len(c.pendingClients), len(c.activeClients))
				//c.broker.RemoveClient(c.conversationID)
			} else if _, ok := c.activeClients[client]; ok {
				delete(c.activeClients, client)
				close(client.send)
				log.Printf("Deregistered a client (%d clients, %d active)", len(c.pendingClients), len(c.activeClients))
				//c.broker.RemoveClient(c.conversationID)
			}

		case message := <-c.broadcast:
			patches, err := dmp.PatchFromText(string(message.content))
			if err != nil {
				return
			}
			*c.doc, _ = dmp.PatchApply(patches, *c.doc)

			for client := range c.activeClients {
				if client != message.sender {
					client.send <- message.content
				}
			}
		}
	}
}
