package mocks

import (
	"local/infra/repo"
	"local/model"

	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of RepositoryInterface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) User() repo.UserRepo {
	args := m.Called()
	return args.Get(0).(repo.UserRepo)
}

func (m *MockRepository) Conversation() repo.ConversationRepo {
	args := m.Called()
	return args.Get(0).(repo.ConversationRepo)
}

func (m *MockRepository) Participant() repo.ParticipantRepo {
	args := m.Called()
	return args.Get(0).(repo.ParticipantRepo)
}

func (m *MockRepository) Message() repo.MessageRepo {
	args := m.Called()
	return args.Get(0).(repo.MessageRepo)
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
