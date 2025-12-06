package conversation

import (
	"fmt"
	"local/infra/repo"
	"local/model"
	"local/service/common"
	"local/util/logger"
	"sort"
)

type ConversationService interface {
	CreateConversation(reqCtx *model.RequestContext, userIDs []uint) model.Response[*model.Conversation]
	GetUserConversations(reqCtx *model.RequestContext, userID uint) model.Response[[]*model.Conversation]
	GetConversationByUserIDs(reqCtx *model.RequestContext, userIDs []uint) model.Response[*model.Conversation]
	GetConversationByID(reqCtx *model.RequestContext, id uint) model.Response[*model.Conversation]
}

type conversationService struct {
	repo repo.RepositoryInterface
}

func convertUserIdsToEntityJoined(userIds []uint) string {
	if len(userIds) != 2 {
		return ""
	}

	entityJoined := ""
	userIdSorted := make([]uint, len(userIds))
	copy(userIdSorted, userIds)
	sort.Slice(userIdSorted, func(i, j int) bool {
		return userIdSorted[i] < userIdSorted[j]
	})

	entityJoined = fmt.Sprintf("user:%d-user:%d", userIdSorted[0], userIdSorted[1])

	return entityJoined
}

func (svc *conversationService) GetConversationByUserIDs(reqCtx *model.RequestContext, userIDs []uint) model.Response[*model.Conversation] {
	logger.Info(reqCtx, "GetConversationByUserIDs called", map[string]interface{}{"user_ids": userIDs})
	if len(userIDs) < 2 {
		return model.BadRequest[*model.Conversation]("At least 2 participants are required")
	}

	response := svc.repo.Conversation().GetByEntityJoined(reqCtx, convertUserIdsToEntityJoined(userIDs))
	return response
}

func (svc *conversationService) CreateConversation(reqCtx *model.RequestContext, userIds []uint) model.Response[*model.Conversation] {
	logger.Info(reqCtx, "CreateConversation called", map[string]interface{}{"user_ids": userIds})
	if len(userIds) < 2 {
		return model.BadRequest[*model.Conversation]("At least 2 participants are required")
	}

	// Create conversation
	conversation := &model.Conversation{
		Type: "private",
		Name: "",
		EntityJoined: convertUserIdsToEntityJoined(userIds),
		UserIds: userIds,
	}

	createResponse := svc.repo.Conversation().Create(reqCtx, conversation)
	if !createResponse.OK() {
		return createResponse
	}

	createdConversation := createResponse.Data

	// Add participants
	for _, userID := range userIds {
		participantResponse := svc.repo.Participant().AddParticipantToConversation(reqCtx, createdConversation.ID, userID)
		if !participantResponse.OK() {
			return model.BadRequest[*model.Conversation]("Failed to add participant to conversation")
		}
	}

	// Return conversation with participants
	queryResponse := svc.repo.Conversation().QueryOne(reqCtx, &model.Conversation{ID: createdConversation.ID})
	return queryResponse
}

func (svc *conversationService) GetUserConversations(reqCtx *model.RequestContext, userID uint) model.Response[[]*model.Conversation] {
	logger.Info(reqCtx, "GetUserConversations called", map[string]interface{}{"user_id": userID})
	response := svc.repo.Conversation().GetByParticipant(reqCtx, userID)
	return response
}

func (svc *conversationService) GetConversationByID(reqCtx *model.RequestContext, id uint) model.Response[*model.Conversation] {
	logger.Info(reqCtx, "GetConversationByID called", map[string]interface{}{"conversation_id": id})
	response := svc.repo.Conversation().QueryOne(reqCtx, &model.Conversation{ID: id})
	return response
}

func NewConversationService(params *common.Params) ConversationService {
	return &conversationService{
		repo: params.Repo,
	}
}
