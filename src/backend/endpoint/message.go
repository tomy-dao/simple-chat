package endpoint

import (
	"context"
	"local/model"
	"local/service/auth"
	"local/service/initial"
	"local/service/message"
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

func (e *MessageEndpoints) CreateMessage(ctx context.Context, request CreateMessageRequest) (Response[*model.Message], error) {
	message, err := e.messageSvc.CreateMessage(ctx, &model.Message{
		ConversationID: request.ConversationID,
		SenderID: request.SenderID,
		Content: request.Content,
		SessionID: request.SessionID,
	})
	if err != nil {
		return Response[*model.Message]{Data: nil, Error: err.Error()}, nil
	}

	return Response[*model.Message]{Data: &message, Error: ""}, nil
}

func (e *MessageEndpoints) GetMessagesByConversationID(ctx context.Context, cvsID uint) (Response[[]*model.Message], error) {
	messages, err := e.messageSvc.GetMessagesByConversationID(ctx, cvsID)
	if err != nil {
		return Response[[]*model.Message]{Data: nil, Error: err.Error()}, nil
	}

	return Response[[]*model.Message]{Data: &messages, Error: ""}, nil
}

func NewMessageEndpoints(params *initial.Service) *MessageEndpoints {
	return &MessageEndpoints{
		messageSvc: params.MessageSvc,
		authSvc: params.AuthSvc,
	}
}