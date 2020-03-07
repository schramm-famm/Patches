package models

type Message struct {
	Type MessageType `json:"type"`
	Data InnerData   `json:"data"`
}

type MessageType int
type UpdateType int

const (
	TypeInit      MessageType = 0
	TypeUpdate    MessageType = 1
	TypeAck       MessageType = 2
	TypeUserJoin  MessageType = 3
	TypeUserLeave MessageType = 4

	UpdateTypeEdit   UpdateType = 0
	UpdateTypeCursor UpdateType = 1
)

type InnerData struct {
	Type        *UpdateType    `json:"type,omitempty"`
	Version     *int           `json:"version,omitempty"`
	Patch       *string        `json:"patch,omitempty"`
	CursorDelta *int           `json:"cursor_delta,omitempty"`
	UserID      *int64         `json:"user_id,omitempty"`
	Content     *string        `json:"content,omitempty"`
	ActiveUsers *map[int64]int `json:"active_users,omitempty"`
}
