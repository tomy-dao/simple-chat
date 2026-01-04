package initial

import (
	"local/service/auth"
	"local/service/common"
	"local/service/conversation"
	"local/service/message"
	"local/service/metrics"
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

	// Initialize Prometheus metrics collector
	metrics.NewPrometheusMetrics(params)

	return Service{
		CvsSvc: CvsSvc,
		AuthSvc: AuthSvc,
		MessageSvc: MessageSvc,
	}
}
