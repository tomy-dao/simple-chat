package integration

import (
	"local/model"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Conversation-specific helpers
func createConversation(t *testing.T, setup *TestSetup, token string, userID uint) *model.Conversation {
	body := map[string]interface{}{"user_id": userID}
	recorder := makeRequest(setup, http.MethodPost, "/api/v1/conversations/", token, body)
	assert.Equal(t, http.StatusOK, recorder.Code)

	response := parseResponse[*model.Conversation](t, recorder)
	assert.True(t, response.OK())
	return response.Data
}

func getConversationByUserID(t *testing.T, setup *TestSetup, token string, userID uint) *model.Conversation {
	path := "/api/v1/conversations/user/" + strconv.Itoa(int(userID))
	recorder := makeRequest(setup, http.MethodGet, path, token, nil)
	assert.Equal(t, http.StatusOK, recorder.Code)

	response := parseResponse[*model.Conversation](t, recorder)
	assert.True(t, response.OK())
	return response.Data
}

func getUserConversations(t *testing.T, setup *TestSetup, token string) []*model.Conversation {
	recorder := makeRequest(setup, http.MethodGet, "/api/v1/conversations/", token, nil)
	assert.Equal(t, http.StatusOK, recorder.Code)

	response := parseResponse[[]*model.Conversation](t, recorder)
	assert.True(t, response.OK())
	return response.Data
}

func TestConversationFlow_CreateAndGet(t *testing.T) {
	setup, err := SetupTestEnvironment()
	assert.NoError(t, err)
	defer setup.CleanupTestEnvironment()

	user1Token := registerAndLogin(t, setup, "user1", "password123")
	user2Token := registerAndLogin(t, setup, "user2", "password123")
	user2ID := getUserID(t, setup, user2Token)

	// Create conversation
	conv := createConversation(t, setup, user1Token, user2ID)
	assert.NotZero(t, conv.ID)

	// Get conversation by user ID
	foundConv := getConversationByUserID(t, setup, user1Token, user2ID)
	assert.Equal(t, conv.ID, foundConv.ID)
}

func TestConversationFlow_GetUserConversations(t *testing.T) {
	setup, err := SetupTestEnvironment()
	assert.NoError(t, err)
	defer setup.CleanupTestEnvironment()

	user1Token := registerAndLogin(t, setup, "user1", "password123")
	user2Token := registerAndLogin(t, setup, "user2", "password123")
	user2ID := getUserID(t, setup, user2Token)

	createConversation(t, setup, user1Token, user2ID)

	// Get all conversations
	conversations := getUserConversations(t, setup, user1Token)
	assert.GreaterOrEqual(t, len(conversations), 1)
}
