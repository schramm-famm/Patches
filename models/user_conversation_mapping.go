package models

// UserConversationMapping represents a user's relationship to a conversation
type UserConversationMapping struct {
	UserID         int64   `json:"user_id,omitempty"`
	ConversationID int64   `json:"conversation_id,omitempty"`
	Role           Role    `json:"role,omitempty"`
	Nickname       *string `json:"nickname,omitempty"`
	Pending        *bool   `json:"pending,omitempty"`
	LastOpened     string  `json:"last_opened,omitempty"`
}

type UserConversationMappingList struct {
	Users []*UserConversationMapping `json:"users"`
}

// Role represents a user's access control rights in a conversation
type Role string

const (
	// Owner is a role that only the original creator of a conversation can have
	// and represents the highest level of privilege
	Owner Role = "owner"

	// Admin is a role that multiple non-creator users in a conversation can
	// have and represents elevated privilege over regular users
	Admin Role = "admin"

	// User is a role that multiple non-creator users in a converation can have
	// and represents the lowest level of privilege
	User Role = "user"
)
