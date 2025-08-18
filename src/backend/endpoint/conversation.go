package endpoint

import (
	"context"
	"local/model"
	"local/service/auth"
	"local/service/conversation"
	"local/service/initial"
)

type ConversationEndpoints struct {
	cvsSvc conversation.ConversationService
	authSvc auth.AuthService
}

type CreateConversationRequest struct {
	UserIDs []uint `json:"user_ids"`
}

type GetUserConversationsRequest struct {
	UserID uint `json:"user_id"`
}


func (e *ConversationEndpoints) CreateConversation(ctx context.Context, userIDs []uint) (Response[*model.Conversation], error) {
	conversation, err := e.cvsSvc.CreateConversation(ctx, userIDs)
	if err != nil {
		return Response[*model.Conversation]{Data: nil, Error: err.Error()}, nil
	}
	
	return Response[*model.Conversation]{Data: &conversation, Error: ""}, nil
}

func (e *ConversationEndpoints) GetConversationByUserIDs(ctx context.Context, userIDs []uint) (*Response[*model.Conversation], error) {
	conversation, err := e.cvsSvc.GetConversationByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return &Response[*model.Conversation]{Data: &conversation, Error: ""}, nil
}

func (e *ConversationEndpoints) GetUserConversations(ctx context.Context, userID uint) (Response[[]*model.Conversation], error) {
	conversations, err := e.cvsSvc.GetUserConversations(ctx, userID)
	if err != nil {
		return Response[[]*model.Conversation]{Data: nil, Error: err.Error()}, nil
	}
	
	return Response[[]*model.Conversation]{Data: &conversations, Error: ""}, nil
}

func NewConversationEndpoints(params *initial.Service) *ConversationEndpoints {
	return &ConversationEndpoints{
		cvsSvc: params.CvsSvc,
		authSvc: params.AuthSvc,
	}
}
