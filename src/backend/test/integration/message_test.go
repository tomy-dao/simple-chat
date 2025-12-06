package integration

import (
	"local/model"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Message-specific helpers
func createMessage(t *testing.T, setup *TestSetup, token string, conversationID uint, content, sessionID string) *model.Message {
	body := map[string]interface{}{
		"content":    content,
		"session_id": sessionID,
	}
	path := "/api/v1/conversations/" + strconv.Itoa(int(conversationID)) + "/messages"
	recorder := makeRequest(setup, http.MethodPost, path, token, body)
	assert.Equal(t, http.StatusOK, recorder.Code)

	response := parseResponse[*model.Message](t, recorder)
	assert.True(t, response.OK())
	return response.Data
}

func getMessages(t *testing.T, setup *TestSetup, token string, conversationID uint) []*model.Message {
	path := "/api/v1/conversations/" + strconv.Itoa(int(conversationID)) + "/messages"
	recorder := makeRequest(setup, http.MethodGet, path, token, nil)
	assert.Equal(t, http.StatusOK, recorder.Code)

	response := parseResponse[[]*model.Message](t, recorder)
	assert.True(t, response.OK())
	return response.Data
}

func TestMessageFlow_CreateAndGet(t *testing.T) {
	setup, err := SetupTestEnvironment()
	assert.NoError(t, err)
	defer setup.CleanupTestEnvironment()

	user1Token := registerAndLogin(t, setup, "user1", "password123")
	user1ID := getUserID(t, setup, user1Token)
	user2Token := registerAndLogin(t, setup, "user2", "password123")
	user2ID := getUserID(t, setup, user2Token)

	// Create conversation
	conv := createConversation(t, setup, user1Token, user2ID)

	// Create message
	msg := createMessage(t, setup, user1Token, conv.ID, "Hello, this is a test message", "test-session-123")
	assert.Equal(t, conv.ID, msg.ConversationID)
	assert.Equal(t, user1ID, msg.SenderID)
	assert.Equal(t, "Hello, this is a test message", msg.Content)

	// Get messages
	messages := getMessages(t, setup, user1Token, conv.ID)
	assert.GreaterOrEqual(t, len(messages), 1)
	assert.Equal(t, "Hello, this is a test message", messages[0].Content)
}

func TestMessageFlow_MultipleMessages(t *testing.T) {
	setup, err := SetupTestEnvironment()
	assert.NoError(t, err)
	defer setup.CleanupTestEnvironment()

	user1Token := registerAndLogin(t, setup, "user1", "password123")
	user2Token := registerAndLogin(t, setup, "user2", "password123")
	user2ID := getUserID(t, setup, user2Token)

	// Create conversation
	conv := createConversation(t, setup, user1Token, user2ID)

	// Create multiple messages
	messages := []string{"Message 1", "Message 2", "Message 3"}
	for i, content := range messages {
		createMessage(t, setup, user1Token, conv.ID, content, "test-session-"+strconv.Itoa(i))
	}

	// Get all messages
	allMessages := getMessages(t, setup, user1Token, conv.ID)
	assert.Equal(t, 3, len(allMessages), "Should have 3 messages")
}
