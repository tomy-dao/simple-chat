package model

import (
	"time"
)

type Message struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ConversationID uint      `json:"conversation_id" gorm:"column:conversation_id;not null"`
	SenderID       uint      `json:"sender_id" gorm:"column:sender_id;not null"`
	Content        string    `json:"content" gorm:"column:content;type:text;not null"`
	MessageType    string    `json:"message_type" gorm:"column:message_type;default:'text'"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	SessionID      string    `json:"session_id,omitempty"`

	Conversation Conversation `json:"conversation" gorm:"foreignKey:ConversationID"`
	Sender       User         `json:"sender" gorm:"foreignKey:SenderID"`
}

func (Message) TableName() string {
	return "messages"
}
