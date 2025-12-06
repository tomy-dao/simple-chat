package repo

import (
	"local/model"
	"local/util/logger"

	"gorm.io/gorm"
)

type ConversationRepo interface {
	QueryOne(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[*model.Conversation]
	QueryMany(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[[]*model.Conversation]
	Create(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[*model.Conversation]
	Update(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[*model.Conversation]
	Delete(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[*model.Conversation]
	GetByParticipant(reqCtx *model.RequestContext, userID uint) model.Response[[]*model.Conversation]
	GetByEntityJoined(reqCtx *model.RequestContext, entityJoined string) model.Response[*model.Conversation]
}

type conversationRepository struct {
	db *gorm.DB
}

func (r *conversationRepository) GetByEntityJoined(reqCtx *model.RequestContext, entityJoined string) model.Response[*model.Conversation] {
	logger.Info(reqCtx, "ConversationRepo.GetByEntityJoined called", map[string]interface{}{"entity_joined": entityJoined})
	var result model.Conversation
	err := r.db.WithContext(reqCtx.Context()).Where("entity_joined = ?", entityJoined).First(&result).Error
	if err != nil {
		return model.NotFound[*model.Conversation]("Conversation not found")
	}
	return model.SuccessResponse(&result, "Conversation retrieved successfully")
}

func (r *conversationRepository) QueryOne(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[*model.Conversation] {
	logger.Info(reqCtx, "ConversationRepo.QueryOne called", map[string]interface{}{"conversation_id": conversation.ID})
	var result model.Conversation
	err := r.db.WithContext(reqCtx.Context()).Preload("Participants.User").Where(conversation).First(&result).Error
	if err != nil {
		return model.NotFound[*model.Conversation]("Conversation not found")
	}
	return model.SuccessResponse(&result, "Conversation retrieved successfully")
}

func (r *conversationRepository) QueryMany(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[[]*model.Conversation] {
	logger.Info(reqCtx, "ConversationRepo.QueryMany called")
	var results []*model.Conversation
	err := r.db.WithContext(reqCtx.Context()).Preload("Participants.User").Where(conversation).Find(&results).Error
	if err != nil {
		return model.InternalError[[]*model.Conversation]("Failed to query conversations")
	}
	return model.SuccessResponse(results, "Conversations retrieved successfully")
}

func (r *conversationRepository) Create(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[*model.Conversation] {
	logger.Info(reqCtx, "ConversationRepo.Create called", map[string]interface{}{"entity_joined": conversation.EntityJoined})
	err := r.db.WithContext(reqCtx.Context()).Create(conversation).Error
	if err != nil {
		return model.BadRequest[*model.Conversation]("Failed to create conversation")
	}
	return model.SuccessResponse(conversation, "Conversation created successfully")
}

func (r *conversationRepository) Update(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[*model.Conversation] {
	logger.Info(reqCtx, "ConversationRepo.Update called", map[string]interface{}{"conversation_id": conversation.ID})
	err := r.db.WithContext(reqCtx.Context()).Save(conversation).Error
	if err != nil {
		return model.BadRequest[*model.Conversation]("Failed to update conversation")
	}
	return model.SuccessResponse(conversation, "Conversation updated successfully")
}

func (r *conversationRepository) Delete(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[*model.Conversation] {
	logger.Info(reqCtx, "ConversationRepo.Delete called", map[string]interface{}{"conversation_id": conversation.ID})
	err := r.db.WithContext(reqCtx.Context()).Delete(conversation).Error
	if err != nil {
		return model.BadRequest[*model.Conversation]("Failed to delete conversation")
	}
	return model.SuccessResponse(conversation, "Conversation deleted successfully")
}

func (r *conversationRepository) GetByParticipant(reqCtx *model.RequestContext, userID uint) model.Response[[]*model.Conversation] {
	logger.Info(reqCtx, "ConversationRepo.GetByParticipant called", map[string]interface{}{"user_id": userID})
	var conversations []*model.Conversation
	err := r.db.WithContext(reqCtx.Context()).
		Preload("Participants.User").
		Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("conversation_participants.user_id = ?", userID).
		Order("last_message_id DESC").
		Find(&conversations).Error
	if err != nil {
		return model.InternalError[[]*model.Conversation]("Failed to get conversations")
	}
	return model.SuccessResponse(conversations, "Conversations retrieved successfully")
}

func NewConversationRepository(db *gorm.DB) (ConversationRepo, error) {
	return &conversationRepository{
		db: db,
	}, nil
}

