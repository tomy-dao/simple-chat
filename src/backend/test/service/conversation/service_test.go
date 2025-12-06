package conversation_test

import (
	"local/model"
	"local/service/common"
	"local/service/conversation"
	"local/test/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConversationService_GetConversationByUserIDs(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockConversationRepo := new(mocks.MockConversationRepo)
	svc := conversation.NewConversationService(&common.Params{Repo: mockRepo})
	reqCtx := &model.RequestContext{}

	t.Run("requires at least two participants", func(t *testing.T) {
		resp := svc.GetConversationByUserIDs(reqCtx, []uint{1})

		assert.Equal(t, model.CodeBadRequest, resp.Code)
		assert.False(t, resp.OK())
		mockRepo.AssertNotCalled(t, "Conversation")
	})

	t.Run("returns conversation when repo succeeds", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockConversationRepo.ExpectedCalls = nil

		mockRepo.On("Conversation").Return(mockConversationRepo)
		mockConversationRepo.On("GetByEntityJoined", reqCtx, "user:2-user:5").
			Return(model.SuccessResponse(&model.Conversation{ID: 10}, "ok"))

		resp := svc.GetConversationByUserIDs(reqCtx, []uint{5, 2})

		assert.True(t, resp.OK())
		assert.Equal(t, uint(10), resp.Data.ID)
		mockRepo.AssertExpectations(t)
		mockConversationRepo.AssertExpectations(t)
	})
}

func TestConversationService_CreateConversation(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockConversationRepo := new(mocks.MockConversationRepo)
	mockParticipantRepo := new(mocks.MockParticipantRepo)
	svc := conversation.NewConversationService(&common.Params{Repo: mockRepo})
	reqCtx := &model.RequestContext{}

	t.Run("validates minimum participants", func(t *testing.T) {
		resp := svc.CreateConversation(reqCtx, []uint{1})
		assert.Equal(t, model.CodeBadRequest, resp.Code)
		mockRepo.AssertNotCalled(t, "Conversation")
	})

	t.Run("returns error when create fails", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockConversationRepo.ExpectedCalls = nil

		mockRepo.On("Conversation").Return(mockConversationRepo)
		mockConversationRepo.On("Create", reqCtx, mock.AnythingOfType("*model.Conversation")).
			Return(model.BadRequest[*model.Conversation]("Failed to create conversation"))

		resp := svc.CreateConversation(reqCtx, []uint{1, 2})

		assert.Equal(t, model.CodeBadRequest, resp.Code)
		mockRepo.AssertExpectations(t)
		mockConversationRepo.AssertExpectations(t)
	})

	t.Run("returns error when adding participant fails", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockConversationRepo.ExpectedCalls = nil
		mockParticipantRepo.ExpectedCalls = nil

		mockRepo.On("Conversation").Return(mockConversationRepo)
		mockRepo.On("Participant").Return(mockParticipantRepo)

		created := &model.Conversation{ID: 7}
		mockConversationRepo.On("Create", reqCtx, mock.AnythingOfType("*model.Conversation")).
			Return(model.SuccessResponse(created, "created"))
		mockParticipantRepo.On("AddParticipantToConversation", reqCtx, uint(7), uint(1)).
			Return(model.BadRequest[*model.ConversationParticipant]("fail"))

		resp := svc.CreateConversation(reqCtx, []uint{1, 3})

		assert.Equal(t, model.CodeBadRequest, resp.Code)
		mockConversationRepo.AssertExpectations(t)
		mockParticipantRepo.AssertExpectations(t)
	})

	t.Run("creates conversation and returns full conversation", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockConversationRepo.ExpectedCalls = nil
		mockParticipantRepo.ExpectedCalls = nil

		mockRepo.On("Conversation").Return(mockConversationRepo)
		mockRepo.On("Participant").Return(mockParticipantRepo)

		created := &model.Conversation{ID: 9}
		mockConversationRepo.On("Create", reqCtx, mock.AnythingOfType("*model.Conversation")).
			Return(model.SuccessResponse(created, "created"))
		mockParticipantRepo.On("AddParticipantToConversation", reqCtx, uint(9), uint(1)).
			Return(model.SuccessResponse(&model.ConversationParticipant{}, "added"))
		mockParticipantRepo.On("AddParticipantToConversation", reqCtx, uint(9), uint(2)).
			Return(model.SuccessResponse(&model.ConversationParticipant{}, "added"))
		mockConversationRepo.On("QueryOne", reqCtx, &model.Conversation{ID: 9}).
			Return(model.SuccessResponse(&model.Conversation{ID: 9}, "ok"))

		resp := svc.CreateConversation(reqCtx, []uint{1, 2})

		assert.True(t, resp.OK())
		assert.Equal(t, uint(9), resp.Data.ID)
		mockConversationRepo.AssertExpectations(t)
		mockParticipantRepo.AssertExpectations(t)
	})
}

func TestConversationService_GetUserConversations(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockConversationRepo := new(mocks.MockConversationRepo)
	svc := conversation.NewConversationService(&common.Params{Repo: mockRepo})
	reqCtx := &model.RequestContext{}

	mockRepo.On("Conversation").Return(mockConversationRepo)
	mockConversationRepo.On("GetByParticipant", reqCtx, uint(5)).
		Return(model.SuccessResponse([]*model.Conversation{{ID: 1}}, "ok"))

	resp := svc.GetUserConversations(reqCtx, 5)

	assert.True(t, resp.OK())
	assert.Len(t, resp.Data, 1)
	mockRepo.AssertExpectations(t)
	mockConversationRepo.AssertExpectations(t)
}

func TestConversationService_GetConversationByID(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockConversationRepo := new(mocks.MockConversationRepo)
	svc := conversation.NewConversationService(&common.Params{Repo: mockRepo})
	reqCtx := &model.RequestContext{}

	mockRepo.On("Conversation").Return(mockConversationRepo)
	mockConversationRepo.On("QueryOne", reqCtx, &model.Conversation{ID: 3}).
		Return(model.SuccessResponse(&model.Conversation{ID: 3}, "ok"))

	resp := svc.GetConversationByID(reqCtx, 3)

	assert.True(t, resp.OK())
	assert.Equal(t, uint(3), resp.Data.ID)
	mockRepo.AssertExpectations(t)
	mockConversationRepo.AssertExpectations(t)
}
