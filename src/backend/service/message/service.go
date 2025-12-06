package message

import (
	"local/client"
	"local/infra/repo"
	"local/model"
	"local/service/auth"
	"local/service/common"
	"local/service/conversation"
	"local/util/logger"
)

type MessageService interface {
	CreateMessage(reqCtx *model.RequestContext, message *model.Message) model.Response[*model.Message]
	GetMessagesByConversationID(reqCtx *model.RequestContext, conversationID uint) model.Response[[]*model.Message]
}

type messageService struct {
	repo repo.RepositoryInterface
	client *client.Client
	authService auth.AuthService
	cvsSvc conversation.ConversationService
}

func (svc *messageService) CreateMessage(reqCtx *model.RequestContext, message *model.Message) model.Response[*model.Message] {
	logger.Info(reqCtx, "CreateMessage called", map[string]interface{}{
		"conversation_id": message.ConversationID,
		"sender_id": message.SenderID,
	})
	createResponse := svc.repo.Message().Create(reqCtx, message)
	if !createResponse.OK() {
		return createResponse
	}

	createdMessage := createResponse.Data

	conversationResponse := svc.cvsSvc.GetConversationByID(reqCtx, message.ConversationID)
	if !conversationResponse.OK() {
		return model.BadRequest[*model.Message]("Conversation not found")
	}

	conversation := conversationResponse.Data
	conversation.LastMessageID = createdMessage.ID
	updateResponse := svc.repo.Conversation().Update(reqCtx, conversation)
	if !updateResponse.OK() {
		return model.BadRequest[*model.Message]("Failed to update conversation")
	}

	createdMessage.Conversation = conversation
	
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
	
	return model.SuccessResponse(createdMessage, "Message created successfully")
}

func (svc *messageService) GetMessagesByConversationID(reqCtx *model.RequestContext, conversationID uint) model.Response[[]*model.Message] {
	logger.Info(reqCtx, "GetMessagesByConversationID called", map[string]interface{}{"conversation_id": conversationID})
	response := svc.repo.Message().GetByConversationID(reqCtx, conversationID)
	return response
}

func NewMessageService(params *common.Params, authService auth.AuthService, cvsSvc conversation.ConversationService) MessageService {
	return &messageService{
		repo: params.Repo,
		client: params.Client,
		authService: authService,
		cvsSvc: cvsSvc,
	}
}
