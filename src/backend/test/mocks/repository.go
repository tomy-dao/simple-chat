package mocks

import (
	"local/infra/repo"
	"local/model"

	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of Repository
type MockRepository struct {
	mock.Mock
	UserRepo        repo.UserRepo
	ConversationRepo repo.ConversationRepo
	ParticipantRepo  repo.ParticipantRepo
	MessageRepo      repo.MessageRepo
}

// MockUserRepo is a mock implementation of UserRepo
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) QueryOne(reqCtx *model.RequestContext, user *model.User) model.Response[*model.User] {
	args := m.Called(reqCtx, user)
	return args.Get(0).(model.Response[*model.User])
}

func (m *MockUserRepo) QueryMany(reqCtx *model.RequestContext, user *model.User) model.Response[[]*model.User] {
	args := m.Called(reqCtx, user)
	return args.Get(0).(model.Response[[]*model.User])
}

func (m *MockUserRepo) Create(reqCtx *model.RequestContext, user *model.User) model.Response[*model.User] {
	args := m.Called(reqCtx, user)
	return args.Get(0).(model.Response[*model.User])
}

func (m *MockUserRepo) Update(reqCtx *model.RequestContext, user *model.User) model.Response[*model.User] {
	args := m.Called(reqCtx, user)
	return args.Get(0).(model.Response[*model.User])
}

func (m *MockUserRepo) Delete(reqCtx *model.RequestContext, user *model.User) model.Response[*model.User] {
	args := m.Called(reqCtx, user)
	return args.Get(0).(model.Response[*model.User])
}

// MockConversationRepo is a mock implementation of ConversationRepo
type MockConversationRepo struct {
	mock.Mock
}

func (m *MockConversationRepo) QueryOne(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[*model.Conversation] {
	args := m.Called(reqCtx, conversation)
	return args.Get(0).(model.Response[*model.Conversation])
}

func (m *MockConversationRepo) QueryMany(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[[]*model.Conversation] {
	args := m.Called(reqCtx, conversation)
	return args.Get(0).(model.Response[[]*model.Conversation])
}

func (m *MockConversationRepo) Create(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[*model.Conversation] {
	args := m.Called(reqCtx, conversation)
	return args.Get(0).(model.Response[*model.Conversation])
}

func (m *MockConversationRepo) Update(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[*model.Conversation] {
	args := m.Called(reqCtx, conversation)
	return args.Get(0).(model.Response[*model.Conversation])
}

func (m *MockConversationRepo) Delete(reqCtx *model.RequestContext, conversation *model.Conversation) model.Response[*model.Conversation] {
	args := m.Called(reqCtx, conversation)
	return args.Get(0).(model.Response[*model.Conversation])
}

func (m *MockConversationRepo) GetByParticipant(reqCtx *model.RequestContext, userID uint) model.Response[[]*model.Conversation] {
	args := m.Called(reqCtx, userID)
	return args.Get(0).(model.Response[[]*model.Conversation])
}

func (m *MockConversationRepo) GetByEntityJoined(reqCtx *model.RequestContext, entityJoined string) model.Response[*model.Conversation] {
	args := m.Called(reqCtx, entityJoined)
	return args.Get(0).(model.Response[*model.Conversation])
}

// MockParticipantRepo is a mock implementation of ParticipantRepo
type MockParticipantRepo struct {
	mock.Mock
}

func (m *MockParticipantRepo) QueryOne(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[*model.ConversationParticipant] {
	args := m.Called(reqCtx, participant)
	return args.Get(0).(model.Response[*model.ConversationParticipant])
}

func (m *MockParticipantRepo) QueryMany(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[[]*model.ConversationParticipant] {
	args := m.Called(reqCtx, participant)
	return args.Get(0).(model.Response[[]*model.ConversationParticipant])
}

func (m *MockParticipantRepo) Create(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[*model.ConversationParticipant] {
	args := m.Called(reqCtx, participant)
	return args.Get(0).(model.Response[*model.ConversationParticipant])
}

func (m *MockParticipantRepo) Update(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[*model.ConversationParticipant] {
	args := m.Called(reqCtx, participant)
	return args.Get(0).(model.Response[*model.ConversationParticipant])
}

func (m *MockParticipantRepo) Delete(reqCtx *model.RequestContext, participant *model.ConversationParticipant) model.Response[*model.ConversationParticipant] {
	args := m.Called(reqCtx, participant)
	return args.Get(0).(model.Response[*model.ConversationParticipant])
}

func (m *MockParticipantRepo) GetByConversationID(reqCtx *model.RequestContext, conversationID uint) model.Response[[]*model.ConversationParticipant] {
	args := m.Called(reqCtx, conversationID)
	return args.Get(0).(model.Response[[]*model.ConversationParticipant])
}

func (m *MockParticipantRepo) GetByUserID(reqCtx *model.RequestContext, userID uint) model.Response[[]*model.ConversationParticipant] {
	args := m.Called(reqCtx, userID)
	return args.Get(0).(model.Response[[]*model.ConversationParticipant])
}

func (m *MockParticipantRepo) GetByConversationAndUser(reqCtx *model.RequestContext, conversationID, userID uint) model.Response[*model.ConversationParticipant] {
	args := m.Called(reqCtx, conversationID, userID)
	return args.Get(0).(model.Response[*model.ConversationParticipant])
}

func (m *MockParticipantRepo) AddParticipantToConversation(reqCtx *model.RequestContext, conversationID, userID uint) model.Response[*model.ConversationParticipant] {
	args := m.Called(reqCtx, conversationID, userID)
	return args.Get(0).(model.Response[*model.ConversationParticipant])
}

func (m *MockParticipantRepo) RemoveParticipantFromConversation(reqCtx *model.RequestContext, conversationID, userID uint) model.Response[string] {
	args := m.Called(reqCtx, conversationID, userID)
	return args.Get(0).(model.Response[string])
}

// MockMessageRepo is a mock implementation of MessageRepo
type MockMessageRepo struct {
	mock.Mock
}

func (m *MockMessageRepo) Create(reqCtx *model.RequestContext, message *model.Message) model.Response[*model.Message] {
	args := m.Called(reqCtx, message)
	return args.Get(0).(model.Response[*model.Message])
}

func (m *MockMessageRepo) GetByConversationID(reqCtx *model.RequestContext, conversationID uint) model.Response[[]*model.Message] {
	args := m.Called(reqCtx, conversationID)
	return args.Get(0).(model.Response[[]*model.Message])
}

// Helper functions to create mocks with default return values
// These can be used to simplify test setup when you need common default behaviors

// NewMockRepositoryWithDefaults creates a MockRepository with default repo mocks
func NewMockRepositoryWithDefaults() (*MockRepository, *MockUserRepo, *MockConversationRepo, *MockParticipantRepo, *MockMessageRepo) {
	mockUserRepo := new(MockUserRepo)
	mockConversationRepo := new(MockConversationRepo)
	mockParticipantRepo := new(MockParticipantRepo)
	mockMessageRepo := new(MockMessageRepo)

	mockRepo := &MockRepository{
		UserRepo:        mockUserRepo,
		ConversationRepo: mockConversationRepo,
		ParticipantRepo:  mockParticipantRepo,
		MessageRepo:      mockMessageRepo,
	}

	return mockRepo, mockUserRepo, mockConversationRepo, mockParticipantRepo, mockMessageRepo
}

// SetupMockUserQueryOneSuccess sets up a successful QueryOne mock for UserRepo
func SetupMockUserQueryOneSuccess(mockUserRepo *MockUserRepo, user *model.User) {
	mockUserRepo.On("QueryOne", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.UserName == user.UserName || u.ID == user.ID
	})).Return(model.SuccessResponse(user, "User found"))
}

// SetupMockUserQueryOneNotFound sets up a not found QueryOne mock for UserRepo
func SetupMockUserQueryOneNotFound(mockUserRepo *MockUserRepo, userName string) {
	mockUserRepo.On("QueryOne", mock.Anything, &model.User{UserName: userName}).
		Return(model.NotFound[*model.User]("User not found"))
}

// SetupMockUserCreateSuccess sets up a successful Create mock for UserRepo
func SetupMockUserCreateSuccess(mockUserRepo *MockUserRepo, user *model.User) {
	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.UserName == user.UserName
	})).Return(model.SuccessResponse(user, "User created successfully"))
}

// SetupMockConversationQueryOneSuccess sets up a successful QueryOne mock for ConversationRepo
func SetupMockConversationQueryOneSuccess(mockConversationRepo *MockConversationRepo, conversation *model.Conversation) {
	mockConversationRepo.On("QueryOne", mock.Anything, &model.Conversation{ID: conversation.ID}).
		Return(model.SuccessResponse(conversation, "Conversation found"))
}

// SetupMockConversationCreateSuccess sets up a successful Create mock for ConversationRepo
func SetupMockConversationCreateSuccess(mockConversationRepo *MockConversationRepo, conversation *model.Conversation) {
	mockConversationRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Conversation")).
		Return(model.SuccessResponse(conversation, "Conversation created successfully"))
}

// SetupMockParticipantAddSuccess sets up a successful AddParticipantToConversation mock
func SetupMockParticipantAddSuccess(mockParticipantRepo *MockParticipantRepo, conversationID, userID uint) {
	mockParticipantRepo.On("AddParticipantToConversation", mock.Anything, conversationID, userID).
		Return(model.SuccessResponse(&model.ConversationParticipant{
			ID:             1,
			ConversationID: conversationID,
			UserID:         userID,
		}, "Participant added successfully"))
}

// SetupMockMessageCreateSuccess sets up a successful Create mock for MessageRepo
func SetupMockMessageCreateSuccess(mockMessageRepo *MockMessageRepo, message *model.Message) {
	mockMessageRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Message")).
		Return(model.SuccessResponse(message, "Message created successfully"))
}

// SetupMockMessageGetByConversationIDSuccess sets up a successful GetByConversationID mock
func SetupMockMessageGetByConversationIDSuccess(mockMessageRepo *MockMessageRepo, conversationID uint, messages []*model.Message) {
	mockMessageRepo.On("GetByConversationID", mock.Anything, conversationID).
		Return(model.SuccessResponse(messages, "Messages retrieved successfully"))
}
