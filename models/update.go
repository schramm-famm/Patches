package models

type Update struct {
	Type        UpdateType `json:"type"`
	Version     int        `json:"version"`
	Patch       string     `json:"patch"`
	CursorDelta int        `json:"cursor_delta"`
}

type UpdateType int

const (
	Edit UpdateType = 0
)
