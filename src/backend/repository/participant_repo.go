package repository

import (
	"context"
	"local/model"

	"gorm.io/gorm"
)

type ParticipantRepo interface {
	QueryOne(ctx context.Context, participant *model.ConversationParticipant) *model.ConversationParticipant
	QueryMany(ctx context.Context, participant *model.ConversationParticipant) []*model.ConversationParticipant
	Create(ctx context.Context, participant *model.ConversationParticipant) *model.ConversationParticipant
	Update(ctx context.Context, participant *model.ConversationParticipant) *model.ConversationParticipant
	Delete(ctx context.Context, participant *model.ConversationParticipant) *model.ConversationParticipant
	GetByConversationID(ctx context.Context, conversationID uint) []*model.ConversationParticipant
	GetByUserID(ctx context.Context, userID uint) []*model.ConversationParticipant
	GetByConversationAndUser(ctx context.Context, conversationID, userID uint) *model.ConversationParticipant
	AddParticipantToConversation(ctx context.Context, conversationID, userID uint) *model.ConversationParticipant
	RemoveParticipantFromConversation(ctx context.Context, conversationID, userID uint) error
}

type participantRepository struct {
	db *gorm.DB
}

func (r *participantRepository) QueryOne(ctx context.Context, participant *model.ConversationParticipant) *model.ConversationParticipant {
	var result model.ConversationParticipant
	err := r.db.WithContext(ctx).Preload("User").Preload("Conversation").Where(participant).First(&result).Error
	if err != nil {
		return nil
	}
	return &result
}

func (r *participantRepository) QueryMany(ctx context.Context, participant *model.ConversationParticipant) []*model.ConversationParticipant {
	var results []*model.ConversationParticipant
	err := r.db.WithContext(ctx).Preload("User").Preload("Conversation").Where(participant).Find(&results).Error
	if err != nil {
		return nil
	}
	return results
}

func (r *participantRepository) Create(ctx context.Context, participant *model.ConversationParticipant) *model.ConversationParticipant {
	err := r.db.WithContext(ctx).Create(participant).Error
	if err != nil {
		return nil
	}
	return participant
}

func (r *participantRepository) Update(ctx context.Context, participant *model.ConversationParticipant) *model.ConversationParticipant {
	err := r.db.WithContext(ctx).Save(participant).Error
	if err != nil {
		return nil
	}
	return participant
}

func (r *participantRepository) Delete(ctx context.Context, participant *model.ConversationParticipant) *model.ConversationParticipant {
	err := r.db.WithContext(ctx).Delete(participant).Error
	if err != nil {
		return nil
	}
	return participant
}

func (r *participantRepository) GetByConversationID(ctx context.Context, conversationID uint) []*model.ConversationParticipant {
	var participants []*model.ConversationParticipant
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Conversation").
		Where("conversation_id = ?", conversationID).
		Find(&participants).Error
	if err != nil {
		return nil
	}
	return participants
}

func (r *participantRepository) GetByUserID(ctx context.Context, userID uint) []*model.ConversationParticipant {
	var participants []*model.ConversationParticipant
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Conversation").
		Where("user_id = ?", userID).
		Find(&participants).Error
	if err != nil {
		return nil
	}
	return participants
}

func (r *participantRepository) GetByConversationAndUser(ctx context.Context, conversationID, userID uint) *model.ConversationParticipant {
	var participant model.ConversationParticipant
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Conversation").
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		First(&participant).Error
	if err != nil {
		return nil
	}
	return &participant
}

func (r *participantRepository) AddParticipantToConversation(ctx context.Context, conversationID, userID uint) *model.ConversationParticipant {
	// Check if participant already exists
	existing := r.GetByConversationAndUser(ctx, conversationID, userID)
	if existing != nil {
		return existing
	}

	participant := &model.ConversationParticipant{
		ConversationID: conversationID,
		UserID:         userID,
	}

	return r.Create(ctx, participant)
}

func (r *participantRepository) RemoveParticipantFromConversation(ctx context.Context, conversationID, userID uint) error {
	err := r.db.WithContext(ctx).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Delete(&model.ConversationParticipant{}).Error
	return err
}

func NewParticipantRepository(db *gorm.DB) (ParticipantRepo, error) {
	return &participantRepository{
		db: db,
	}, nil
}
