package endpoint

import "local/service/initial"


type Endpoints struct {
	Auth *AuthEndpoints
	Conversation *ConversationEndpoints
	Message *MessageEndpoints
}

func NewEndpoints(params *initial.Service) *Endpoints {
	auth := NewAuthEndpoints(params)
	conversation := NewConversationEndpoints(params)
	message := NewMessageEndpoints(params)
	return &Endpoints{
		Auth: auth,
		Conversation: conversation,
		Message: message,
	}
}