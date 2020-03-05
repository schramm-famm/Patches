package websockets

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"patches/models"
	"patches/utils"

	gorillaws "github.com/gorilla/websocket"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Conversation manages all WebSocket connections in a single conversation.
type Conversation struct {
	conversationID int64
	doc            string
	clients        map[*Client]bool
	version        int

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
		version:        0,
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

func (c *Conversation) processMessage(msg *models.Message) (bool, error) {
	if msg.Type != models.TypeUpdate {
		errMsg := fmt.Sprintf("Message is not of type %d", models.TypeUpdate)
		return false, errors.New(errMsg)
	}

	update := msg.Data
	if *update.Version < 1 {
		errMsg := fmt.Sprintf("Update has invalid version number %d", update.Version)
		return false, errors.New(errMsg)
	}

	patches, err := dmp.PatchFromText(*update.Patch)
	if err != nil {
		return false, err
	}
	if len(patches) != 1 {
		errMsg := "Update must contain one patch"
		return false, errors.New(errMsg)
	}

	newDoc, okList := dmp.PatchApply(patches, c.doc)
	if !okList[0] {
		return false, nil
	}
	c.doc = newDoc

	if *update.Version != c.version+1 {
		*update.Version = c.version + 1
	}

	return true, nil
}

// Run waits on a Conversation's three channels for clients to be added, clients
// to be removed, and messages to be broadcast. Only one of these operations may
// be performed at a time.
func (c *Conversation) Run() {
	for {
		select {
		case client := <-c.register:
			init := models.Message{
				Type: models.TypeInit,
				Data: models.InnerData{
					Version: utils.IntPtr(c.version),
					Content: &c.doc,
				},
			}
			initMessage, err := json.Marshal(init)
			if err != nil {
				log.Print("Failed to encode initial conversation data: ", err)
				close(client.send)
				continue
			}

			err = client.conn.WriteMessage(gorillaws.TextMessage, initMessage)
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

			msg := models.Message{}
			if err := json.Unmarshal(message.content, &msg); err != nil {
				log.Printf("Failed to parse WebSocket message content: %v", err)
				c.deleteClient(message.sender)
				continue
			}

			originalVersion := *msg.Data.Version
			ok, err := c.processMessage(&msg)
			if err != nil {
				log.Printf("Failed to apply update: %v", err)
				c.deleteClient(message.sender)
				continue
			}
			if !ok {
				log.Printf("Patch %s could not be applied", *msg.Data.Patch)
				continue
			}

			broadcastMessageBytes := message.content
			if *msg.Data.Version != originalVersion {
				broadcastMessageBytes, err = json.Marshal(msg)
				if err != nil {
					log.Printf("Failed to encode update as byte array: %v", err)
					c.deleteClient(message.sender)
					continue
				}
			}

			for client := range c.clients {
				if client != message.sender {
					client.send <- broadcastMessageBytes
				}
			}
			c.version++
			message.sender.position += *msg.Data.CursorDelta

			ackMessage := models.Message{
				Type: models.TypeAck,
				Data: models.InnerData{
					Version: &originalVersion,
				},
			}
			ackMessageBytes, err := json.Marshal(ackMessage)
			if err != nil {
				log.Printf("Failed to encode update as byte array: %v", err)
				c.deleteClient(message.sender)
				continue
			}
			message.sender.send <- ackMessageBytes
		}
	}
}
