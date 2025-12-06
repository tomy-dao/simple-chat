package message_test

import (
	"local/client"
	"local/model"
	"local/service/auth"
	"local/service/common"
	"local/service/conversation"
	"local/service/message"
	"local/test/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConversationService struct {
	mock.Mock
}

func (m *MockConversationService) CreateConversation(reqCtx *model.RequestContext, userIDs []uint) model.Response[*model.Conversation] {
	args := m.Called(reqCtx, userIDs)
	return args.Get(0).(model.Response[*model.Conversation])
}

func (m *MockConversationService) GetUserConversations(reqCtx *model.RequestContext, userID uint) model.Response[[]*model.Conversation] {
	args := m.Called(reqCtx, userID)
	return args.Get(0).(model.Response[[]*model.Conversation])
}

func (m *MockConversationService) GetConversationByUserIDs(reqCtx *model.RequestContext, userIDs []uint) model.Response[*model.Conversation] {
	args := m.Called(reqCtx, userIDs)
	return args.Get(0).(model.Response[*model.Conversation])
}

func (m *MockConversationService) GetConversationByID(reqCtx *model.RequestContext, id uint) model.Response[*model.Conversation] {
	args := m.Called(reqCtx, id)
	return args.Get(0).(model.Response[*model.Conversation])
}

type MockAuthService struct{}

func (m *MockAuthService) Authenticate(reqCtx *model.RequestContext) model.Response[uint] {
	return model.Response[uint]{}
}
func (m *MockAuthService) ParseToken(tokenStr string) model.Response[*auth.JWTClaims] {
	return model.Response[*auth.JWTClaims]{}
}
func (m *MockAuthService) CheckToken(reqCtx *model.RequestContext, token string) model.Response[bool] {
	return model.Response[bool]{}
}
func (m *MockAuthService) GetMe(reqCtx *model.RequestContext) model.Response[*model.User] {
	return model.Response[*model.User]{}
}
func (m *MockAuthService) Register(reqCtx *model.RequestContext, userName, password string) model.Response[*model.User] {
	return model.Response[*model.User]{}
}
func (m *MockAuthService) Login(reqCtx *model.RequestContext, userName, password string) model.Response[string] {
	return model.Response[string]{}
}
func (m *MockAuthService) Logout(reqCtx *model.RequestContext, token string) model.Response[string] {
	return model.Response[string]{}
}
func (m *MockAuthService) GetUsers(reqCtx *model.RequestContext) model.Response[[]*model.User] {
	return model.Response[[]*model.User]{}
}

type MockSocketClient struct {
	mock.Mock
}

func (m *MockSocketClient) Broadcast(message *model.BroadcastMessage) {
	m.Called(message)
}

func newMessageService(mockRepo *mocks.MockRepository, mockSocket *MockSocketClient, cvsSvc conversation.ConversationService) message.MessageService {
	params := &common.Params{
		Repo:   mockRepo,
		Client: &client.Client{SocketClient: mockSocket},
	}
	return message.NewMessageService(params, &MockAuthService{}, cvsSvc)
}

func TestMessageService_CreateMessage_CreateFails(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockMessageRepo := new(mocks.MockMessageRepo)
	mockConversationRepo := new(mocks.MockConversationRepo)
	mockConversationService := new(MockConversationService)
	mockSocket := new(MockSocketClient)

	svc := newMessageService(mockRepo, mockSocket, mockConversationService)
	reqCtx := &model.RequestContext{}
	msg := &model.Message{ConversationID: 1, SenderID: 2, Content: "hello", SessionID: "s1"}

	mockRepo.On("Message").Return(mockMessageRepo)
	mockRepo.On("Conversation").Return(mockConversationRepo)
	mockMessageRepo.On("Create", reqCtx, msg).
		Return(model.BadRequest[*model.Message]("create fail"))

	resp := svc.CreateMessage(reqCtx, msg)

	assert.False(t, resp.OK())
	assert.Equal(t, model.CodeBadRequest, resp.Code)
	mockMessageRepo.AssertExpectations(t)
	mockConversationService.AssertNotCalled(t, "GetConversationByID", mock.Anything, mock.Anything)
	mockSocket.AssertNotCalled(t, "Broadcast", mock.Anything)
}

func TestMessageService_CreateMessage_ConversationMissing(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockMessageRepo := new(mocks.MockMessageRepo)
	mockConversationRepo := new(mocks.MockConversationRepo)
	mockConversationService := new(MockConversationService)
	mockSocket := new(MockSocketClient)

	svc := newMessageService(mockRepo, mockSocket, mockConversationService)
	reqCtx := &model.RequestContext{}
	msg := &model.Message{ConversationID: 2, SenderID: 3, Content: "hello"}

	created := &model.Message{ID: 10, ConversationID: 2, SenderID: 3, Content: "hello"}

	mockRepo.On("Message").Return(mockMessageRepo)
	mockRepo.On("Conversation").Return(mockConversationRepo)
	mockMessageRepo.On("Create", reqCtx, msg).
		Return(model.SuccessResponse(created, "created"))
	mockConversationService.On("GetConversationByID", reqCtx, uint(2)).
		Return(model.NotFound[*model.Conversation]("Conversation not found"))

	resp := svc.CreateMessage(reqCtx, msg)

	assert.False(t, resp.OK())
	assert.Equal(t, model.CodeBadRequest, resp.Code)
	assert.Contains(t, resp.ErrorString(), "Conversation not found")
	mockMessageRepo.AssertExpectations(t)
	mockConversationService.AssertExpectations(t)
	mockSocket.AssertNotCalled(t, "Broadcast", mock.Anything)
}

func TestMessageService_CreateMessage_UpdateFails(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockMessageRepo := new(mocks.MockMessageRepo)
	mockConversationRepo := new(mocks.MockConversationRepo)
	mockConversationService := new(MockConversationService)
	mockSocket := new(MockSocketClient)

	svc := newMessageService(mockRepo, mockSocket, mockConversationService)
	reqCtx := &model.RequestContext{}
	msg := &model.Message{ConversationID: 3, SenderID: 4, Content: "hi"}

	created := &model.Message{ID: 11, ConversationID: 3, SenderID: 4, Content: "hi"}
	conversationData := &model.Conversation{ID: 3}

	mockRepo.On("Message").Return(mockMessageRepo)
	mockRepo.On("Conversation").Return(mockConversationRepo)

	mockMessageRepo.On("Create", reqCtx, msg).
		Return(model.SuccessResponse(created, "created"))
	mockConversationService.On("GetConversationByID", reqCtx, uint(3)).
		Return(model.SuccessResponse(conversationData, "ok"))
	mockConversationRepo.On("Update", reqCtx, mock.MatchedBy(func(c *model.Conversation) bool {
		return c.LastMessageID == created.ID
	})).Return(model.BadRequest[*model.Conversation]("Failed to update conversation"))

	resp := svc.CreateMessage(reqCtx, msg)

	assert.False(t, resp.OK())
	assert.Equal(t, model.CodeBadRequest, resp.Code)
	assert.Contains(t, resp.ErrorString(), "Failed to update conversation")
	mockMessageRepo.AssertExpectations(t)
	mockConversationService.AssertExpectations(t)
	mockConversationRepo.AssertExpectations(t)
	mockSocket.AssertNotCalled(t, "Broadcast", mock.Anything)
}

func TestMessageService_CreateMessage_Success(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockMessageRepo := new(mocks.MockMessageRepo)
	mockConversationRepo := new(mocks.MockConversationRepo)
	mockConversationService := new(MockConversationService)
	mockSocket := new(MockSocketClient)

	svc := newMessageService(mockRepo, mockSocket, mockConversationService)
	reqCtx := &model.RequestContext{}
	msg := &model.Message{ConversationID: 4, SenderID: 5, Content: "content", SessionID: "session-1"}

	created := &model.Message{ID: 20, ConversationID: 4, SenderID: 5, Content: "content", SessionID: "session-1"}
	conversationData := &model.Conversation{
		ID: 4,
		Participants: []*model.ConversationParticipant{
			{UserID: 1},
			{UserID: 5},
		},
	}

	mockRepo.On("Message").Return(mockMessageRepo)
	mockRepo.On("Conversation").Return(mockConversationRepo)

	mockMessageRepo.On("Create", reqCtx, msg).
		Return(model.SuccessResponse(created, "created"))
	mockConversationService.On("GetConversationByID", reqCtx, uint(4)).
		Return(model.SuccessResponse(conversationData, "ok"))
	mockConversationRepo.On("Update", reqCtx, mock.MatchedBy(func(c *model.Conversation) bool {
		return c.LastMessageID == created.ID
	})).Return(model.SuccessResponse(conversationData, "updated"))
	mockSocket.On("Broadcast", mock.MatchedBy(func(b *model.BroadcastMessage) bool {
		return b.Event == "message" &&
			len(b.UserIds) == 2 &&
			b.UserIds[0] == 1 &&
			b.UserIds[1] == 5 &&
			b.SessionId == msg.SessionID
	})).Return()

	resp := svc.CreateMessage(reqCtx, msg)

	assert.True(t, resp.OK())
	assert.NotNil(t, resp.Data)
	assert.Equal(t, uint(20), resp.Data.ID)
	assert.NotNil(t, resp.Data.Conversation)
	assert.Equal(t, uint(20), resp.Data.Conversation.LastMessageID)
	mockMessageRepo.AssertExpectations(t)
	mockConversationService.AssertExpectations(t)
	mockConversationRepo.AssertExpectations(t)
	mockSocket.AssertExpectations(t)
}

func TestMessageService_GetMessagesByConversationID(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockMessageRepo := new(mocks.MockMessageRepo)
	mockConversationRepo := new(mocks.MockConversationRepo)
	mockSocket := new(MockSocketClient)

	svc := newMessageService(mockRepo, mockSocket, &MockConversationService{})
	reqCtx := &model.RequestContext{}

	mockRepo.On("Message").Return(mockMessageRepo)
	mockRepo.On("Conversation").Return(mockConversationRepo)
	mockMessageRepo.On("GetByConversationID", reqCtx, uint(9)).
		Return(model.SuccessResponse([]*model.Message{{ID: 1}}, "ok"))

	resp := svc.GetMessagesByConversationID(reqCtx, 9)

	assert.True(t, resp.OK())
	assert.Len(t, resp.Data, 1)
	mockMessageRepo.AssertExpectations(t)
}
