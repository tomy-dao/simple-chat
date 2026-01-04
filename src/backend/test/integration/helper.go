package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"local/client"
	"local/config"
	"local/endpoint"
	"local/infra/repo"
	"local/model"
	"local/service/common"
	"local/service/initial"
	httpTransport "local/transport/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestSetup contains all dependencies for integration tests
type TestSetup struct {
	DB        *gorm.DB
	Repo      *repo.Repository
	Endpoints *endpoint.Endpoints
	Router    *gin.Engine
}

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate all models
	err = db.AutoMigrate(
		&model.User{},
		&model.Conversation{},
		&model.ConversationParticipant{},
		&model.Message{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// SetupTestEnvironment creates a complete test environment
func SetupTestEnvironment() (*TestSetup, error) {
	// Setup test database
	db, err := SetupTestDB()
	if err != nil {
		return nil, err
	}

	// Create test repository using NewRepositoryWithDB (works for both test and production)
	testRepo, err := repo.NewRepositoryWithDB(db)
	if err != nil {
		return nil, err
	}

	// Setup config for testing
	config.Config.JwtSecret = "test-jwt-secret-key-for-integration-tests"

	// Setup rate limiting config for tests (with permissive defaults)
	// Individual tests can override these values
	if config.Config.RateLimitRequestsPerMin == 0 {
		config.Config.RateLimitRequestsPerMin = 1000 // High default to not interfere with tests
	}
	if config.Config.RateLimitBurst == 0 {
		config.Config.RateLimitBurst = 100 // High default to not interfere with tests
	}

	// Create services
	initParams := &model.InitParams{
		ServiceName: "test-service",
		Ctx:         context.Background(),
	}
	clt := client.NewClient(initParams)

	svc := initial.NewService(&common.Params{
		Repo:   testRepo,
		Client: clt,
	})

	// Create endpoints
	endpoints := endpoint.NewEndpoints(&svc)

	// Create HTTP router
	router := httpTransport.MakeHttpTransport(initParams, endpoints)

	return &TestSetup{
		DB:        db,
		Repo:      testRepo,
		Endpoints: endpoints,
		Router:    router,
	}, nil
}

// CleanupTestEnvironment cleans up test data
func (ts *TestSetup) CleanupTestEnvironment() error {
	if ts.DB != nil {
		// Get all table names from database
		var tables []string
		ts.DB.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&tables)
		
		// Delete all data from each table
		for _, table := range tables {
			ts.DB.Exec("DELETE FROM " + table)
		}
	}
	return nil
}

// HTTP helpers - Common utilities for all tests
func makeRequest(setup *TestSetup, method, path, token string, body interface{}) *httptest.ResponseRecorder {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	setup.Router.ServeHTTP(recorder, req)
	return recorder
}

func parseResponse[T any](t *testing.T, recorder *httptest.ResponseRecorder) model.Response[T] {
	var response model.Response[T]
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	return response
}
