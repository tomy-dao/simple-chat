package repo

import (
	"local/model"
	"local/util/logger"

	"gorm.io/gorm"
)

type ParticipantRepo interface {
	QueryOne(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[*model.ConversationParticipant]
	QueryMany(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[[]*model.ConversationParticipant]
	Create(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[*model.ConversationParticipant]
	Update(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[*model.ConversationParticipant]
	Delete(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[*model.ConversationParticipant]
	GetByConversationID(reqCtx *model.RequestContext, conversationID uint) model.Response[[]*model.ConversationParticipant]
	GetByUserID(reqCtx *model.RequestContext, userID uint) model.Response[[]*model.ConversationParticipant]
	GetByConversationAndUser(reqCtx *model.RequestContext, conversationID, userID uint) model.Response[*model.ConversationParticipant]
	AddParticipantToConversation(reqCtx *model.RequestContext, conversationID, userID uint) model.Response[*model.ConversationParticipant]
	RemoveParticipantFromConversation(reqCtx *model.RequestContext, conversationID, userID uint) model.Response[string]
}

type participantRepository struct {
	db *gorm.DB
}

func (r *participantRepository) QueryOne(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[*model.ConversationParticipant] {
	logger.Info(reqCtx, "ParticipantRepo.QueryOne called", map[string]interface{}{
		"conversation_id": participant.ConversationID,
		"user_id": participant.UserID,
	})
	var result model.ConversationParticipant
	err := r.db.WithContext(reqCtx.Context()).Preload("User").Preload("Conversation").Where(participant).First(&result).Error
	if err != nil {
		return model.NotFound[*model.ConversationParticipant]("Participant not found")
	}
	return model.SuccessResponse(&result, "Participant retrieved successfully")
}

func (r *participantRepository) QueryMany(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[[]*model.ConversationParticipant] {
	logger.Info(reqCtx, "ParticipantRepo.QueryMany called")
	var results []*model.ConversationParticipant
	err := r.db.WithContext(reqCtx.Context()).Preload("User").Preload("Conversation").Where(participant).Find(&results).Error
	if err != nil {
		return model.InternalError[[]*model.ConversationParticipant]("Failed to query participants")
	}
	return model.SuccessResponse(results, "Participants retrieved successfully")
}

func (r *participantRepository) Create(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[*model.ConversationParticipant] {
	logger.Info(reqCtx, "ParticipantRepo.Create called", map[string]interface{}{
		"conversation_id": participant.ConversationID,
		"user_id": participant.UserID,
	})
	err := r.db.WithContext(reqCtx.Context()).Create(participant).Error
	if err != nil {
		return model.BadRequest[*model.ConversationParticipant]("Failed to create participant")
	}
	return model.SuccessResponse(participant, "Participant created successfully")
}

func (r *participantRepository) Update(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[*model.ConversationParticipant] {
	logger.Info(reqCtx, "ParticipantRepo.Update called", map[string]interface{}{"participant_id": participant.ID})
	err := r.db.WithContext(reqCtx.Context()).Save(participant).Error
	if err != nil {
		return model.BadRequest[*model.ConversationParticipant]("Failed to update participant")
	}
	return model.SuccessResponse(participant, "Participant updated successfully")
}

func (r *participantRepository) Delete(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[*model.ConversationParticipant] {
	logger.Info(reqCtx, "ParticipantRepo.Delete called", map[string]interface{}{"participant_id": participant.ID})
	err := r.db.WithContext(reqCtx.Context()).Delete(participant).Error
	if err != nil {
		return model.BadRequest[*model.ConversationParticipant]("Failed to delete participant")
	}
	return model.SuccessResponse(participant, "Participant deleted successfully")
}

func (r *participantRepository) GetByConversationID(reqCtx *model.RequestContext, conversationID uint) model.Response[[]*model.ConversationParticipant] {
	logger.Info(reqCtx, "ParticipantRepo.GetByConversationID called", map[string]interface{}{"conversation_id": conversationID})
	var participants []*model.ConversationParticipant
	err := r.db.WithContext(reqCtx.Context()).
		Preload("User").
		Preload("Conversation").
		Where("conversation_id = ?", conversationID).
		Find(&participants).Error
	if err != nil {
		return model.InternalError[[]*model.ConversationParticipant]("Failed to get participants")
	}
	return model.SuccessResponse(participants, "Participants retrieved successfully")
}

func (r *participantRepository) GetByUserID(reqCtx *model.RequestContext, userID uint) model.Response[[]*model.ConversationParticipant] {
	logger.Info(reqCtx, "ParticipantRepo.GetByUserID called", map[string]interface{}{"user_id": userID})
	var participants []*model.ConversationParticipant
	err := r.db.WithContext(reqCtx.Context()).
		Preload("User").
		Preload("Conversation").
		Where("user_id = ?", userID).
		Find(&participants).Error
	if err != nil {
		return model.InternalError[[]*model.ConversationParticipant]("Failed to get participants")
	}
	return model.SuccessResponse(participants, "Participants retrieved successfully")
}

func (r *participantRepository) GetByConversationAndUser(reqCtx *model.RequestContext, conversationID, userID uint) model.Response[*model.ConversationParticipant] {
	logger.Info(reqCtx, "ParticipantRepo.GetByConversationAndUser called", map[string]interface{}{
		"conversation_id": conversationID,
		"user_id": userID,
	})
	var participant model.ConversationParticipant
	err := r.db.WithContext(reqCtx.Context()).
		Preload("User").
		Preload("Conversation").
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		First(&participant).Error
	if err != nil {
		return model.NotFound[*model.ConversationParticipant]("Participant not found")
	}
	return model.SuccessResponse(&participant, "Participant retrieved successfully")
}

func (r *participantRepository) AddParticipantToConversation(reqCtx *model.RequestContext, conversationID, userID uint) model.Response[*model.ConversationParticipant] {
	logger.Info(reqCtx, "ParticipantRepo.AddParticipantToConversation called", map[string]interface{}{
		"conversation_id": conversationID,
		"user_id": userID,
	})
	// Check if participant already exists
	existing := r.GetByConversationAndUser(reqCtx, conversationID, userID)
	if existing.OK() {
		return existing
	}

	participant := &model.ConversationParticipant{
		ConversationID: conversationID,
		UserID:         userID,
	}

	return r.Create(reqCtx, participant)
}

func (r *participantRepository) RemoveParticipantFromConversation(reqCtx *model.RequestContext, conversationID, userID uint) model.Response[string] {
	logger.Info(reqCtx, "ParticipantRepo.RemoveParticipantFromConversation called", map[string]interface{}{
		"conversation_id": conversationID,
		"user_id": userID,
	})
	err := r.db.WithContext(reqCtx.Context()).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Delete(&model.ConversationParticipant{}).Error
	if err != nil {
		return model.BadRequest[string]("Failed to remove participant")
	}
	return model.SuccessResponse("", "Participant removed successfully")
}

func NewParticipantRepository(db *gorm.DB) (ParticipantRepo, error) {
	return &participantRepository{
		db: db,
	}, nil
}

