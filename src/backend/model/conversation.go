package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Conversation struct {
	ID          		uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Type        		string    `json:"type" gorm:"column:type;not null"`
	Name        		string    `json:"name" gorm:"column:name"`
	CreatedAt   		time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   		time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	EntityJoined 		string `json:"entity_joined" gorm:"type:text"`
	UserIds 				UserIds `json:"user_ids" gorm:"type:json"`
	
	Participants 		[]ConversationParticipant `json:"participants" gorm:"foreignKey:ConversationID"`
}

// UserIds is a custom type to handle JSON serialization for MySQL
type UserIds []uint

// Value implements the driver.Valuer interface
func (u UserIds) Value() (driver.Value, error) {
	if u == nil {
		return nil, nil
	}
	return json.Marshal(u)
}

// Scan implements the sql.Scanner interface
func (u *UserIds) Scan(value interface{}) error {
	if value == nil {
		*u = nil
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, u)
	case string:
		return json.Unmarshal([]byte(v), u)
	default:
		return nil
	}
}

type ConversationParticipant struct {
	ID             uint `json:"id" gorm:"primaryKey;autoIncrement"`
	ConversationID uint `json:"conversation_id" gorm:"column:conversation_id;not null"`
	UserID         uint `json:"user_id" gorm:"column:user_id;not null"`
	JoinedAt       time.Time `json:"joined_at" gorm:"column:joined_at;autoCreateTime"`
	
	Conversation Conversation `json:"conversation" gorm:"foreignKey:ConversationID"`
	User         User         `json:"user" gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Conversation) TableName() string {
	return "conversations"
}

func (ConversationParticipant) TableName() string {
	return "conversation_participants"
}
