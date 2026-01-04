package repo

import (
	"local/model"
	"local/util/logger"

	"gorm.io/gorm"
)

type MessageRepo interface {
	Create(reqCtx *model.RequestContext, message *model.Message) model.Response[*model.Message]
	GetByConversationID(reqCtx *model.RequestContext, conversationID uint) model.Response[[]*model.Message]
	Count(reqCtx *model.RequestContext) (int64, error)
}

type messageRepository struct {
	db *gorm.DB
}

func (r *messageRepository) Create(reqCtx *model.RequestContext, message *model.Message) model.Response[*model.Message] {
	logger.Info(reqCtx, "MessageRepo.Create called", map[string]interface{}{
		"conversation_id": message.ConversationID,
		"sender_id": message.SenderID,
	})
	err := r.db.WithContext(reqCtx.Context()).Create(message).Error
	if err != nil {
		return model.BadRequest[*model.Message]("Failed to create message")
	}
	return model.SuccessResponse(message, "Message created successfully")
}

func (r *messageRepository) GetByConversationID(reqCtx *model.RequestContext, conversationID uint) model.Response[[]*model.Message] {
	logger.Info(reqCtx, "MessageRepo.GetByConversationID called", map[string]interface{}{"conversation_id": conversationID})
	var messages []*model.Message
	err := r.db.WithContext(reqCtx.Context()).Where("conversation_id = ?", conversationID).Order("id DESC").Find(&messages).Error
	if err != nil {
		return model.InternalError[[]*model.Message]("Failed to get messages")
	}
	return model.SuccessResponse(messages, "Messages retrieved successfully")
}

func (r *messageRepository) Count(reqCtx *model.RequestContext) (int64, error) {
	logger.Info(reqCtx, "MessageRepo.Count called")
	var count int64
	err := r.db.WithContext(reqCtx.Context()).Model(&model.Message{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func NewMessageRepository(db *gorm.DB) (MessageRepo, error) {
	return &messageRepository{db: db}, nil
}

