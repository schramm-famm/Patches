package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"patches/kafka"
	"patches/models"
	"patches/protocol"
	"time"

	gorillaws "github.com/gorilla/websocket"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Conversation manages all WebSocket connections in a single conversation.
type Conversation struct {
	conversationID int64
	doc            string
	clients        map[*Client]bool
	version        int
	checkpoint     map[int]*Checkpoint

	register   chan *Client
	unregister chan *Client
	broadcast  chan *BroadcastMessage
	errc       chan error

	db          models.Datastore
	kafkaWriter *kafka.Writer
}

// BroadcastMessage stores the content and sender of a WebSocket message that is
// meant to be broadcasted to all other clients in a conversation.
type BroadcastMessage struct {
	content []byte
	sender  *Client
}

// Checkpoint stores active users' carets for a version, the caret position of
// the sender and the delta of the patch that brought the conversation to this
// version, and the outstanding Sync's for the version.
type Checkpoint struct {
	activeUsers map[int64]protocol.Caret
	senderCaret protocol.Caret
	delta       protocol.Delta
	syncsLeft   map[int64]bool
}

var dmp *diffmatchpatch.DiffMatchPatch = diffmatchpatch.New()

// NewConversation creates a new Conversation struct.
func NewConversation(
	conversationID int64,
	doc string,
	db models.Datastore,
	kafkaWriter *kafka.Writer,
) *Conversation {
	return &Conversation{
		conversationID: conversationID,
		doc:            doc,
		clients:        make(map[*Client]bool),
		version:        0,
		checkpoint: map[int]*Checkpoint{
			0: &Checkpoint{
				activeUsers: make(map[int64]protocol.Caret),
				syncsLeft:   make(map[int64]bool),
			},
		},
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan *BroadcastMessage),
		errc:        make(chan error),
		db:          db,
		kafkaWriter: kafkaWriter,
	}
}

// sendMessage sends a message to a single receiving client.
func (c *Conversation) sendMessage(msg protocol.Message, receiver *Client) error {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	receiver.send <- messageBytes

	return nil
}

// broadcastMessage sends a message to all clients except a specified sender.
func (c *Conversation) broadcastMessage(msg protocol.Message, sender *Client) error {
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
func (c *Conversation) handleEditUpdate(msg protocol.Message, sender *Client) error {
	update := msg.Data

	if update.Type == nil || update.Version == nil || update.Patch == nil || update.Delta == nil {
		return fmt.Errorf(`update (EDIT) is missing required fields in "data"`)
	}

	if update.Delta.CaretStart == nil || update.Delta.CaretEnd == nil || update.Delta.Doc == nil {
		return fmt.Errorf(`update (EDIT) is missing required fields in "data.delta"`)
	}

	if *update.Version < 1 {
		return fmt.Errorf("update (EDIT) has invalid version number %d", update.Version)
	}

	patches, err := dmp.PatchFromText(*update.Patch)
	if err != nil {
		return err
	}
	if len(patches) != 1 {
		return fmt.Errorf("update (EDIT) must contain one patch")
	}

	newDoc, okList := dmp.PatchApply(patches, c.doc)
	if !okList[0] {
		log.Println("Received invalid patch. Will not be broadcasted or acknowledged.")
		return nil
	}

	if *msg.Data.Version != c.version+1 {
		*msg.Data.Version = c.version + 1
	}

	// Broadcast Update (EDIT) message to all existing clients
	msg.Data.UserID = &sender.userID
	if err := c.broadcastMessage(msg, sender); err != nil {
		return err
	}

	// Publish Update (EDIT) message to Kafka topic
	go func() {
		if err := c.kafkaWriter.PublishUpdate(msg, c.conversationID); err != nil {
			c.errc <- err
		}
	}()

	// TODO: Write relevant Update (EDIT) message data to TimescaleDB hypertable

	ackMessage := protocol.Message{
		Type: protocol.TypeAck,
		Data: protocol.InnerData{
			Version: msg.Data.Version,
		},
	}
	if err := c.sendMessage(ackMessage, sender); err != nil {
		return err
	}

	newCheckpoint := &Checkpoint{
		activeUsers: make(map[int64]protocol.Caret),
		senderCaret: sender.caret,
		delta:       *update.Delta,
		syncsLeft:   make(map[int64]bool),
	}

	// Update all other clients' carets
	for client := range c.clients {
		if client != sender {
			client.caret = client.caret.ShiftCaret(sender.caret, *update.Delta)
			newCheckpoint.syncsLeft[client.userID] = true
			newCheckpoint.activeUsers[client.userID] = client.caret
		}
	}

	// Update the sender's caret
	sender.caret.Start += *update.Delta.CaretStart
	sender.caret.End += *update.Delta.CaretEnd
	newCheckpoint.activeUsers[sender.userID] = sender.caret

	c.checkpoint[*msg.Data.Version] = newCheckpoint

	if len(newCheckpoint.syncsLeft) == 0 {
		syncMessage := protocol.Message{
			Type: protocol.TypeSync,
			Data: protocol.InnerData{
				Version: msg.Data.Version,
			},
		}
		if err := c.handleSync(syncMessage, sender); err != nil {
			return err
		}
	}

	c.version++
	c.doc = newDoc

	return nil
}

// handleCursorUpdate processes an Update message of subtype Cursor and
// broadcasts it out to all clients in the conversation that aren't the sender.
func (c *Conversation) handleCursorUpdate(msg protocol.Message, sender *Client) error {
	update := msg.Data

	if update.Type == nil || update.Delta == nil || update.Version == nil {
		return fmt.Errorf(`update (CURSOR) is missing required fields in "data"`)
	}

	if update.Delta.CaretStart == nil || update.Delta.CaretEnd == nil {
		return fmt.Errorf(`update (CURSOR) is missing required fields in "data.delta"`)
	}

	// Broadcast Update (CURSOR) message to all existing clients
	msg.Data.UserID = &sender.userID
	if err := c.broadcastMessage(msg, sender); err != nil {
		return err
	}

	// Apply delta to the sender's caret position at its version checkpoint
	updateCheckpoint, ok := c.checkpoint[*update.Version]
	if !ok {
		return fmt.Errorf("Version %d does not exist in checkpoints", *update.Version)
	}
	senderCaret := updateCheckpoint.activeUsers[sender.userID]
	senderCaret.Start += *update.Delta.CaretStart
	senderCaret.End += *update.Delta.CaretEnd
	updateCheckpoint.activeUsers[sender.userID] = senderCaret
	sender.caret = senderCaret

	// Adjust sender's caret position at subsequent version checkpoints
	for v := *update.Version + 1; v <= c.version; v++ {
		prevCaret := c.checkpoint[v-1].activeUsers[sender.userID]
		c.checkpoint[v].activeUsers[sender.userID] = prevCaret.ShiftCaret(
			c.checkpoint[v].senderCaret,
			c.checkpoint[v].delta,
		)
	}

	return nil
}

// handleSync processes Sync messages and tracks outstanding syncs for
// checkpoint versions. A version checkpoint is removed when there are no
// outstanding syncs.
func (c *Conversation) handleSync(msg protocol.Message, sender *Client) error {
	sync := msg.Data

	if sync.Version == nil {
		return fmt.Errorf(`sync is missing required fields in "data"`)
	}

	if _, ok := c.checkpoint[*sync.Version]; !ok {
		log.Printf("Checkpoint for version %d does not exist. Removing all previous checkpoints", *sync.Version)
		for i := *sync.Version - 1; ; i-- {
			if _, ok := c.checkpoint[i]; ok {
				delete(c.checkpoint, i)
			} else {
				break
			}
		}
		return nil
	}

	delete(c.checkpoint[*sync.Version].syncsLeft, sender.userID)

	if len(c.checkpoint[*sync.Version].syncsLeft) == 0 {
		delete(c.checkpoint, *sync.Version-1)
		if err := c.broadcastMessage(msg, nil); err != nil {
			return err
		}
	}

	return nil
}

// registerClient starts tracking a client in the conversation, sends the client
// an Init message, and broadcasts a UserJoin message to the rest of the
// clients.
func (c *Conversation) registerClient(client *Client) error {
	// Create and send Init message to the new client
	init := protocol.Message{
		Type: protocol.TypeInit,
		Data: protocol.InnerData{
			Version: &c.version,
			Content: &c.doc,
		},
	}
	if len(c.clients) > 0 {
		activeUsers := make(map[int64]protocol.Caret)
		for client := range c.clients {
			activeUsers[client.userID] = protocol.Caret{
				Start: client.caret.Start,
				End:   client.caret.End,
			}
		}
		init.Data.ActiveUsers = &activeUsers
	}
	initMessage, err := json.Marshal(init)
	if err != nil {
		return err
	}
	client.conn.SetWriteDeadline(time.Now().Add(writeWait))
	err = client.conn.WriteMessage(gorillaws.TextMessage, initMessage)
	if err != nil {
		return err
	}

	// Create and broadcast UserJoin message to all existing clients
	userJoinMsg := protocol.Message{
		Type: protocol.TypeUserJoin,
		Data: protocol.InnerData{
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

	for version := range c.checkpoint {
		delete(c.checkpoint[version].activeUsers, client.userID)
		delete(c.checkpoint[version].syncsLeft, client.userID)
	}

	delete(c.clients, client)
	close(client.send)
	log.Printf("Unregistered a client in conversation %d (%d active)", c.conversationID, len(c.clients))

	// Create and broadcast UserLeave message to all existing clients
	userLeaveMsg := protocol.Message{
		Type: protocol.TypeUserLeave,
		Data: protocol.InnerData{
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

	msg := protocol.Message{}
	if err := json.Unmarshal(broadcastMsg.content, &msg); err != nil {
		return fmt.Errorf("failed to parse WebSocket message content: %v", err)
	}

	if msg.Type != protocol.TypeUpdate && msg.Type != protocol.TypeSync {
		return fmt.Errorf(
			"message is not of type %d or type %d",
			protocol.TypeUpdate,
			protocol.TypeSync,
		)
	}

	if msg.Type == protocol.TypeSync {
		return c.handleSync(msg, broadcastMsg.sender)
	}

	if msg.Data.Type == nil {
		return fmt.Errorf(`update is missing required "type" field in "data"`)
	}

	switch *msg.Data.Type {
	case protocol.UpdateTypeEdit:
		if err := c.handleEditUpdate(msg, broadcastMsg.sender); err != nil {
			return err
		}

	case protocol.UpdateTypeCursor:
		if err := c.handleCursorUpdate(msg, broadcastMsg.sender); err != nil {
			return err
		}

	default:
		return fmt.Errorf("update has invalid subtype %d", *msg.Data.Type)
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
				log.Printf("Shutting down conversation %d", c.conversationID)
				return
			}
			if err := c.processBroadcast(broadcastMsg); err != nil {
				log.Print("Failed to process broadcast message: ", err)
				c.unregisterClient(broadcastMsg.sender)
			}

		case err := <-c.errc:
			log.Print("Error occured during asynchronous action: ", err)
			for client := range c.clients {
				c.unregisterClient(client)
			}

		}
	}
}
