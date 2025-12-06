package endpoint

import (
	"local/model"
	"local/service/auth"
	"local/service/conversation"
	"local/service/initial"
	"local/util/logger"
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


func (e *ConversationEndpoints) CreateConversation(reqCtx *model.RequestContext, userIDs []uint) model.Response[*model.Conversation] {
	logger.Info(reqCtx, "ConversationEndpoints.CreateConversation called", map[string]interface{}{"user_ids": userIDs})
	return e.cvsSvc.CreateConversation(reqCtx, userIDs)
}

func (e *ConversationEndpoints) GetConversationByUserIDs(reqCtx *model.RequestContext, userIDs []uint) model.Response[*model.Conversation] {
	logger.Info(reqCtx, "ConversationEndpoints.GetConversationByUserIDs called", map[string]interface{}{"user_ids": userIDs})
	return e.cvsSvc.GetConversationByUserIDs(reqCtx, userIDs)
}

func (e *ConversationEndpoints) GetUserConversations(reqCtx *model.RequestContext, userID uint) model.Response[[]*model.Conversation] {
	logger.Info(reqCtx, "ConversationEndpoints.GetUserConversations called", map[string]interface{}{"user_id": userID})
	return e.cvsSvc.GetUserConversations(reqCtx, userID)
}

func NewConversationEndpoints(params *initial.Service) *ConversationEndpoints {
	return &ConversationEndpoints{
		cvsSvc: params.CvsSvc,
		authSvc: params.AuthSvc,
	}
}
