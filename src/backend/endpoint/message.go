package endpoint

import (
	"local/model"
	"local/service/auth"
	"local/service/initial"
	"local/service/message"
	"local/util/logger"
)

type MessageEndpoints struct {
	authSvc auth.AuthService
	messageSvc message.MessageService
}

type CreateMessageRequest struct {
	ConversationID uint `json:"conversation_id"`
	SenderID uint `json:"sender_id"`
	Content string `json:"content"`
	SessionID string `json:"session_id"`
}

func (e *MessageEndpoints) CreateMessage(reqCtx *model.RequestContext, request CreateMessageRequest) model.Response[*model.Message] {
	logger.Info(reqCtx, "MessageEndpoints.CreateMessage called", map[string]interface{}{
		"conversation_id": request.ConversationID,
		"sender_id": request.SenderID,
	})
	return e.messageSvc.CreateMessage(reqCtx, &model.Message{
		ConversationID: request.ConversationID,
		SenderID: request.SenderID,
		Content: request.Content,
		SessionID: request.SessionID,
	})
}

func (e *MessageEndpoints) GetMessagesByConversationID(reqCtx *model.RequestContext, cvsID uint) model.Response[[]*model.Message] {
	logger.Info(reqCtx, "MessageEndpoints.GetMessagesByConversationID called", map[string]interface{}{"conversation_id": cvsID})
	return e.messageSvc.GetMessagesByConversationID(reqCtx, cvsID)
}

func NewMessageEndpoints(params *initial.Service) *MessageEndpoints {
	return &MessageEndpoints{
		messageSvc: params.MessageSvc,
		authSvc: params.AuthSvc,
	}
}