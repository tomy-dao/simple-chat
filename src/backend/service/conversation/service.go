package conversation

import (
	"context"
	"errors"
	"fmt"
	"local/model"
	"local/repository"
	"local/service/common"
	"sort"
)

type ConversationService interface {
	CreateConversation(ctx context.Context, userIDs []uint) (*model.Conversation, error)
	GetUserConversations(ctx context.Context, userID uint) ([]*model.Conversation, error)
	GetConversationByUserIDs(ctx context.Context, userIDs []uint) (*model.Conversation, error)
	GetConversationByID(ctx context.Context, id uint) (*model.Conversation, error)
}

type conversationService struct {
	repo repository.RepositoryInterface
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

func (svc *conversationService) GetConversationByUserIDs(ctx context.Context, userIDs []uint) (*model.Conversation, error) {
	if len(userIDs) < 2 {
		return nil, errors.New("at least 2 participants are required")
	}

	conversation := svc.repo.Conversation().GetByEntityJoined(ctx, convertUserIdsToEntityJoined(userIDs))
	if conversation == nil {
		return nil, errors.New("conversation not found")
	}

	return conversation, nil
}

func (svc *conversationService) CreateConversation(ctx context.Context, userIds []uint) (*model.Conversation, error) {
	if len(userIds) < 2 {
		return nil, errors.New("at least 2 participants are required")
	}

	// Create conversation
	conversation := &model.Conversation{
		Type: "private",
		Name: "",
		EntityJoined: convertUserIdsToEntityJoined(userIds),
		UserIds: userIds,
	}

	createdConversation := svc.repo.Conversation().Create(ctx, conversation)
	if createdConversation == nil {
		return nil, errors.New("failed to create conversation")
	}

	// Add participants
	for _, userID := range userIds {
		participant := svc.repo.Participant().AddParticipantToConversation(ctx, createdConversation.ID, userID)
		if participant == nil {
			return nil, errors.New("failed to add participant to conversation")
		}
	}

	// Return conversation with participants
	return svc.repo.Conversation().QueryOne(ctx, &model.Conversation{ID: createdConversation.ID}), nil
}

func (svc *conversationService) GetUserConversations(ctx context.Context, userID uint) ([]*model.Conversation, error) {
	conversations := svc.repo.Conversation().GetByParticipant(ctx, userID)
	if conversations == nil {
		return []*model.Conversation{}, nil
	}

	return conversations, nil
}

func (svc *conversationService) GetConversationByID(ctx context.Context, id uint) (*model.Conversation, error) {
	conversation := svc.repo.Conversation().QueryOne(ctx, &model.Conversation{ID: id})
	if conversation == nil {
		return nil, errors.New("conversation not found")
	}
	return conversation, nil
}

func NewConversationService(params *common.Params) ConversationService {
	return &conversationService{
		repo: params.Repo,
	}
}
