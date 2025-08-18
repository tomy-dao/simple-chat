package repository

import (
	"context"
	"local/model"

	"gorm.io/gorm"
)

type ConversationRepo interface {
	QueryOne(ctx context.Context, conversation *model.Conversation) *model.Conversation
	QueryMany(ctx context.Context, conversation *model.Conversation) []*model.Conversation
	Create(ctx context.Context, conversation *model.Conversation) *model.Conversation
	Update(ctx context.Context, conversation *model.Conversation) *model.Conversation
	Delete(ctx context.Context, conversation *model.Conversation) *model.Conversation
	GetByParticipant(ctx context.Context, userID uint) []*model.Conversation
	GetByEntityJoined(ctx context.Context, entityJoined string) *model.Conversation
}

type conversationRepository struct {
	db *gorm.DB
}

func (r *conversationRepository) GetByEntityJoined(ctx context.Context, entityJoined string) *model.Conversation {
	var result model.Conversation
	err := r.db.WithContext(ctx).Where("entity_joined = ?", entityJoined).First(&result).Error
	if err != nil {
		return nil
	}
	return &result
}

func (r *conversationRepository) QueryOne(ctx context.Context, conversation *model.Conversation) *model.Conversation {
	var result model.Conversation
	err := r.db.WithContext(ctx).Preload("Participants.User").Where(conversation).First(&result).Error
	if err != nil {
		return nil
	}
	return &result
}

func (r *conversationRepository) QueryMany(ctx context.Context, conversation *model.Conversation) []*model.Conversation {
	var results []*model.Conversation
	err := r.db.WithContext(ctx).Preload("Participants.User").Where(conversation).Find(&results).Error
	if err != nil {
		return nil
	}
	return results
}

func (r *conversationRepository) Create(ctx context.Context, conversation *model.Conversation) *model.Conversation {
	err := r.db.WithContext(ctx).Create(conversation).Error
	if err != nil {
		return nil
	}
	return conversation
}

func (r *conversationRepository) Update(ctx context.Context, conversation *model.Conversation) *model.Conversation {
	err := r.db.WithContext(ctx).Save(conversation).Error
	if err != nil {
		return nil
	}
	return conversation
}

func (r *conversationRepository) Delete(ctx context.Context, conversation *model.Conversation) *model.Conversation {
	err := r.db.WithContext(ctx).Delete(conversation).Error
	if err != nil {
		return nil
	}
	return conversation
}

func (r *conversationRepository) GetByParticipant(ctx context.Context, userID uint) []*model.Conversation {
	var conversations []*model.Conversation
	err := r.db.WithContext(ctx).
		Preload("Participants.User").
		Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("conversation_participants.user_id = ?", userID).
		Find(&conversations).Error
	if err != nil {
		return nil
	}
	return conversations
}

func NewConversationRepository(db *gorm.DB) (ConversationRepo, error) {
	return &conversationRepository{
		db: db,
	}, nil
}

