package models

type Message struct {
	Type MessageType `json:"type"`
	Data InnerData   `json:"data"`
}

type MessageType int
type UpdateType int

const (
	TypeInit   MessageType = 0
	TypeUpdate MessageType = 1
	TypeAck    MessageType = 2

	TypeEdit UpdateType = 0
)

type InnerData struct {
	Type        *UpdateType `json:"type,omitempty"`
	Version     *int        `json:"version,omitempty"`
	Patch       *string     `json:"patch,omitempty"`
	CursorDelta *int        `json:"cursor_delta,omitempty"`
	Content     *string     `json:"content,omitempty"`
}
