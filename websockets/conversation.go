package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"patches/models"

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
	broadcast  chan *BroadcastMessage
}

// BroadcastMessage stores the content and sender of a WebSocket message.
type BroadcastMessage struct {
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
		broadcast:      make(chan *BroadcastMessage),
	}
}

// sendMessage sends a message to a single receiving client
func (c *Conversation) sendMessage(msg models.Message, receiver *Client) error {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	receiver.send <- messageBytes

	return nil
}

// broadcastMessage sends a message to all clients except a specified sender
func (c *Conversation) broadcastMessage(msg models.Message, sender *Client) error {
	broadcastMessageBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for client := range c.clients {
		if client != sender {
			client.send <- broadcastMessageBytes
		}
	}

	return nil
}

// handleEditUpdate processes an Update message of subtype Edit and broadcasts
// it out to all clients in the conversation that aren't the sender.
func (c *Conversation) handleEditUpdate(msg models.Message, sender *Client) error {
	update := msg.Data

	if update.Type == nil || update.Version == nil || update.Patch == nil || update.CursorDelta == nil {
		return fmt.Errorf(`Update (EDIT) is missing required fields in "data"`)
	}

	if *update.Version < 1 {
		return fmt.Errorf("Update has invalid version number %d", update.Version)
	}

	patches, err := dmp.PatchFromText(*update.Patch)
	if err != nil {
		return err
	}
	if len(patches) != 1 {
		return fmt.Errorf("Update must contain one patch")
	}

	newDoc, okList := dmp.PatchApply(patches, c.doc)
	if !okList[0] {
		log.Println("Received invalid patch. Will not be broadcasted or acknowledged.")
		return nil
	}
	c.doc = newDoc

	if *update.Version != c.version+1 {
		*msg.Data.Version = c.version + 1
	}

	msg.Data.UserID = &sender.userID
	if err := c.broadcastMessage(msg, sender); err != nil {
		return err
	}

	sender.position += *msg.Data.CursorDelta
	c.version++

	ackMessage := models.Message{
		Type: models.TypeAck,
		Data: models.InnerData{
			Version: update.Version,
		},
	}
	if err := c.sendMessage(ackMessage, sender); err != nil {
		return err
	}

	return nil
}

// handleCursorUpdate processes an Update message of subtype Cursor and
// broadcasts it out to all clients in the conversation that aren't the sender.
func (c *Conversation) handleCursorUpdate(msg models.Message, sender *Client) error {
	update := msg.Data

	if update.Type == nil || update.CursorDelta == nil {
		return fmt.Errorf(`Update (CURSOR) is missing required fields in "data"`)
	}

	msg.Data.UserID = &sender.userID
	if err := c.broadcastMessage(msg, sender); err != nil {
		return err
	}
	sender.position += *msg.Data.CursorDelta

	return nil
}

// registerClient starts tracking a client in the conversation, sends the client
// an Init message, and broadcasts a UserJoin message to the rest of the
// clients.
func (c *Conversation) registerClient(client *Client) error {
	// Create and send Init message to the new client
	init := models.Message{
		Type: models.TypeInit,
		Data: models.InnerData{
			Version: &c.version,
			Content: &c.doc,
		},
	}
	if len(c.clients) > 0 {
		activeUsers := make(map[int64]int)
		for client := range c.clients {
			activeUsers[client.userID] = client.position
		}
		init.Data.ActiveUsers = &activeUsers
	}
	initMessage, err := json.Marshal(init)
	if err != nil {
		return err
	}
	err = client.conn.WriteMessage(gorillaws.TextMessage, initMessage)
	if err != nil {
		return err
	}

	// Create and broadcast UserJoin message to all existing clients
	userJoinMsg := models.Message{
		Type: models.TypeUserJoin,
		Data: models.InnerData{
			UserID: &client.userID,
		},
	}
	if err := c.broadcastMessage(userJoinMsg, nil); err != nil {
		return err
	}

	c.clients[client] = true
	log.Printf("Registered a client in conversation %d (%d active)", c.conversationID, len(c.clients))

	return nil
}

// unregisterClient stops tracking a client in the conversation and broadcasts
// a UserLeave message to the rest of the clients.
func (c *Conversation) unregisterClient(client *Client) error {
	if _, ok := c.clients[client]; !ok {
		log.Printf("Attempted to unregister an inactive client in conversation %d", c.conversationID)
		return nil
	}

	delete(c.clients, client)
	close(client.send)
	log.Printf("Unregistered a client in conversation %d (%d active)", c.conversationID, len(c.clients))

	// Create and broadcast UserLeave message to all existing clients
	userLeaveMsg := models.Message{
		Type: models.TypeUserLeave,
		Data: models.InnerData{
			UserID: &client.userID,
		},
	}
	if err := c.broadcastMessage(userLeaveMsg, nil); err != nil {
		return err
	}

	return nil
}

// processBroadcast processes a received Update message and handles it according
// to the message's subtype.
func (c *Conversation) processBroadcast(broadcastMsg *BroadcastMessage) error {
	if _, ok := c.clients[broadcastMsg.sender]; !ok {
		log.Printf("Attempted to broadcast from an inactive client in conversation %d", c.conversationID)
		return nil
	}

	msg := models.Message{}
	if err := json.Unmarshal(broadcastMsg.content, &msg); err != nil {
		return fmt.Errorf("failed to parse WebSocket message content: %v", err)
	}

	if msg.Type != models.TypeUpdate {
		return fmt.Errorf("message is not of type %d", models.TypeUpdate)
	}

	if msg.Data.Type == nil {
		return fmt.Errorf(`update is missing required "type" field in "data"`)
	}

	switch *msg.Data.Type {
	case models.UpdateTypeEdit:
		if err := c.handleEditUpdate(msg, broadcastMsg.sender); err != nil {
			return err
		}

	case models.UpdateTypeCursor:
		if err := c.handleCursorUpdate(msg, broadcastMsg.sender); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Update has invalid subtype %d", *msg.Data.Type)
	}

	return nil
}

// Run waits on a Conversation's three channels for clients to be added, clients
// to be removed, and messages to be broadcast. Only one of these operations may
// be performed at a time.
func (c *Conversation) Run() {
	for {
		select {
		case client := <-c.register:
			if err := c.registerClient(client); err != nil {
				log.Print("Error occured while registering new client: ", err)
			}

		case client := <-c.unregister:
			if err := c.unregisterClient(client); err != nil {
				log.Print("Error occured while unregistering client: ", err)
			}

		case broadcastMsg, ok := <-c.broadcast:
			if !ok {
				close(c.register)
				close(c.unregister)
				log.Printf("Shutting down conversation %d", c.conversationID)
				return
			}
			if err := c.processBroadcast(broadcastMsg); err != nil {
				log.Print("Failed to process broadcast message: ", err)
				c.unregisterClient(broadcastMsg.sender)
			}

		}
	}
}
