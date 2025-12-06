package integration

import (
	"encoding/json"
	"local/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthFlow_RegisterAndLogin(t *testing.T) {
	setup, err := SetupTestEnvironment()
	assert.NoError(t, err)
	defer setup.CleanupTestEnvironment()

	// Register
	registerUser(t, setup, "testuser", "password123")

	// Login
	token := loginUser(t, setup, "testuser", "password123")
	assert.NotEmpty(t, token)

	// GetMe
	user := getUserID(t, setup, token)
	assert.NotZero(t, user)
}

func TestAuthFlow_InvalidCredentials(t *testing.T) {
	setup, err := SetupTestEnvironment()
	assert.NoError(t, err)
	defer setup.CleanupTestEnvironment()

	registerUser(t, setup, "testuser", "password123")

	// Try login with wrong password
	body := map[string]string{"username": "testuser", "password": "wrongpassword"}
	recorder := makeRequest(setup, http.MethodPost, "/api/v1/login", "", body)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	response := parseResponse[interface{}](t, recorder)
	assert.False(t, response.OK())
	assert.Equal(t, model.CodeUnauthorized, response.Code)
}

func TestAuthFlow_DuplicateRegistration(t *testing.T) {
	setup, err := SetupTestEnvironment()
	assert.NoError(t, err)
	defer setup.CleanupTestEnvironment()

	registerUser(t, setup, "testuser", "password123")

	// Try to register same username again
	body := map[string]string{"username": "testuser", "password": "password123"}
	recorder := makeRequest(setup, http.MethodPost, "/api/v1/register", "", body)

	assert.Equal(t, http.StatusConflict, recorder.Code)
	response := parseResponse[interface{}](t, recorder)
	assert.False(t, response.OK())
	assert.Equal(t, model.CodeConflict, response.Code)
}

// Auth helpers - Used across multiple test files
func registerUser(t *testing.T, setup *TestSetup, username, password string) {
	body := map[string]string{"username": username, "password": password}
	recorder := makeRequest(setup, http.MethodPost, "/api/v1/register", "", body)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func loginUser(t *testing.T, setup *TestSetup, username, password string) string {
	body := map[string]string{"username": username, "password": password}
	recorder := makeRequest(setup, http.MethodPost, "/api/v1/login", "", body)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	return response.Data.Token
}

func registerAndLogin(t *testing.T, setup *TestSetup, username, password string) string {
	registerUser(t, setup, username, password)
	return loginUser(t, setup, username, password)
}

func getUserID(t *testing.T, setup *TestSetup, token string) uint {
	recorder := makeRequest(setup, http.MethodGet, "/api/v1/me", token, nil)
	assert.Equal(t, http.StatusOK, recorder.Code)

	response := parseResponse[*model.User](t, recorder)
	assert.True(t, response.OK())
	return response.Data.ID
}
