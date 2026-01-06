package model

// SessionNotice stores notification messages in Session, often deleted after use
type SessionNotice struct {
	Type    string // Message type
	Content string // Message content
}
