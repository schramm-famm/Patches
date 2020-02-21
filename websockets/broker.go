package websockets

import (
	"net/http"
	"patches/models"
	"sync"

	gorillaws "github.com/gorilla/websocket"
)

type Broker struct {
	sync.Mutex
	internal   map[int64]*Conversation
	db         models.Datastore
	httpClient *http.Client
}

func NewBroker(db models.Datastore, httpClient *http.Client) *Broker {
	return &Broker{
		internal:   make(map[int64]*Conversation),
		db:         db,
		httpClient: httpClient,
	}
}

func (b *Broker) CreateClient(conversationID int64, conn *gorillaws.Conn) *Client {
	b.Lock()
	defer b.Unlock()

	conversation := b.internal[conversationID]
	if conversation == nil {
		conversation = NewConversation(conversationID, b)
		b.internal[conversationID] = conversation
		go conversation.Run()
	}

	client := NewClient(-1, conn, conversation)
	conversation.register <- client
	return client
}

func (b *Broker) removeConversation(conversationID int64) {
	b.Lock()
	defer b.Unlock()

	if conversation := b.internal[conversationID]; conversation != nil {
		b.internal[conversationID] = nil
	}
}
