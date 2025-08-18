package message

import (
	"context"
	"errors"
	"local/client"
	"local/model"
	"local/repository"
	"local/service/auth"
	"local/service/common"
	"local/service/conversation"
)

type MessageService interface {
	CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error)
	GetMessagesByConversationID(ctx context.Context, conversationID uint) ([]*model.Message, error)
}

type messageService struct {
	repo repository.RepositoryInterface
	client *client.Client
	authService auth.AuthService
	cvsSvc conversation.ConversationService
}

func (svc *messageService) CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
	createdMessage := svc.repo.Message().Create(ctx, message)
	if createdMessage == nil {
		return nil, errors.New("failed to create message")
	}

	conversation, _ := svc.cvsSvc.GetConversationByID(ctx, message.ConversationID)

	userIds := []int{}
	for _, participant := range conversation.Participants {
		userIds = append(userIds, int(participant.UserID))
	}

	svc.client.SocketClient.Broadcast(&model.BroadcastMessage{
		UserIds: userIds,
		SessionId: message.SessionID,
		Event: "message",
		Payload: map[string]interface{}{
			"message": message,
		},
	})
	
	return createdMessage, nil
}

func (svc *messageService) GetMessagesByConversationID(ctx context.Context, conversationID uint) ([]*model.Message, error) {
	messages := svc.repo.Message().GetByConversationID(ctx, conversationID)
	if messages == nil {
		return nil, errors.New("failed to get messages")
	}

	return messages, nil
}

func NewMessageService(params *common.Params, authService auth.AuthService, cvsSvc conversation.ConversationService) MessageService {
	return &messageService{
		repo: params.Repo,
		client: params.Client,
		authService: authService,
		cvsSvc: cvsSvc,
	}
}
