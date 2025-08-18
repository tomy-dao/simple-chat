package initial

import (
	"local/service/auth"
	"local/service/common"
	"local/service/conversation"
	"local/service/message"
)

type Service struct {
	CvsSvc conversation.ConversationService
	AuthSvc auth.AuthService
	MessageSvc message.MessageService
}


func NewService(params *common.Params) Service {
	CvsSvc := conversation.NewConversationService(params)
	AuthSvc := auth.NewAuthService(params)
	MessageSvc := message.NewMessageService(params, AuthSvc, CvsSvc)
	
	return Service{
		CvsSvc: CvsSvc,
		AuthSvc: AuthSvc,
		MessageSvc: MessageSvc,
	}
}
