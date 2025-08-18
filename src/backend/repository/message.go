package repository

import (
	"context"
	"local/model"

	"gorm.io/gorm"
)

type MessageRepo interface {
	Create(ctx context.Context, message *model.Message) *model.Message
	GetByConversationID(ctx context.Context, conversationID uint) []*model.Message
}

type messageRepository struct {
	db *gorm.DB
}

func (r *messageRepository) Create(ctx context.Context, message *model.Message) *model.Message {
	err := r.db.WithContext(ctx).Create(message).Error
	if err != nil {
		return nil
	}
	return message
}

func (r *messageRepository) GetByConversationID(ctx context.Context, conversationID uint) []*model.Message {
	var messages []*model.Message
	err := r.db.WithContext(ctx).Where("conversation_id = ?", conversationID).Order("id DESC").Find(&messages).Error
	if err != nil {
		return nil
	}
	return messages
}

func NewMessageRepository(db *gorm.DB) (MessageRepo, error) {
	return &messageRepository{db: db}, nil
}