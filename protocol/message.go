package protocol

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
	TypeSync      MessageType = 3
	TypeUserJoin  MessageType = 4
	TypeUserLeave MessageType = 5

	UpdateTypeEdit   UpdateType = 0
	UpdateTypeCursor UpdateType = 1
)

type InnerData struct {
	Type        *UpdateType      `json:"type,omitempty"`
	Version     *int             `json:"version,omitempty"`
	Patch       *string          `json:"patch,omitempty"`
	Delta       *Delta           `json:"delta,omitempty"`
	UserID      *int64           `json:"user_id,omitempty"`
	Content     *string          `json:"content,omitempty"`
	ActiveUsers *map[int64]Caret `json:"active_users,omitempty"`
}

type Delta struct {
	CaretStart *int `json:"caret_start,omitempty"`
	CaretEnd   *int `json:"caret_end,omitempty"`
	Doc        *int `json:"doc,omitempty"`
}
