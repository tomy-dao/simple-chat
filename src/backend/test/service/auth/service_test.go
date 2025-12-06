package auth_test

import (
	"local/model"
	"local/service/auth"
	"local/test/mocks"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)


func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name           string
		userName       string
		password       string
		setupMocks     func(*mocks.MockRepository, *mocks.MockUserRepo)
		expectedCode   int
		expectedOK     bool
		expectedError  string
	}{
		{
			name:     "successful registration",
			userName: "testuser",
			password: "password123",
			setupMocks: func(mockRepo *mocks.MockRepository, mockUserRepo *mocks.MockUserRepo) {
				// User doesn't exist
				mockUserRepo.On("QueryOne", mock.Anything, &model.User{UserName: "testuser"}).
					Return(model.NotFound[*model.User]("User not found"))
				// Create user
				mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
					return u.UserName == "testuser" && u.Password != ""
				})).Return(model.SuccessResponse(&model.User{
					ID:       1,
					UserName: "testuser",
					Password: "hashed",
				}, "User created successfully"))
				mockRepo.On("User").Return(mockUserRepo)
			},
			expectedCode: model.CodeSuccess,
			expectedOK:   true,
		},
		{
			name:     "empty username",
			userName: "",
			password: "password123",
			setupMocks: func(mockRepo *mocks.MockRepository, mockUserRepo *mocks.MockUserRepo) {
				// No repo calls expected, service returns early
			},
			expectedCode:  model.CodeBadRequest,
			expectedOK:    false,
			expectedError: "Username and password are required",
		},
		{
			name:     "empty password",
			userName: "testuser",
			password: "",
			setupMocks: func(mockRepo *mocks.MockRepository, mockUserRepo *mocks.MockUserRepo) {
				// No repo calls expected, service returns early
			},
			expectedCode:  model.CodeBadRequest,
			expectedOK:    false,
			expectedError: "Username and password are required",
		},
		{
			name:     "user already exists",
			userName: "existinguser",
			password: "password123",
			setupMocks: func(mockRepo *mocks.MockRepository, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.On("QueryOne", mock.Anything, &model.User{UserName: "existinguser"}).
					Return(model.SuccessResponse(&model.User{
						ID:       1,
						UserName: "existinguser",
					}, "User found"))
				mockRepo.On("User").Return(mockUserRepo)
			},
			expectedCode:  model.CodeConflict,
			expectedOK:    false,
			expectedError: "User already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockRepository)
			mockUserRepo := new(mocks.MockUserRepo)
			tt.setupMocks(mockRepo, mockUserRepo)

			svc := auth.NewTestAuthService(mockRepo, "test-secret-key-for-testing-only")

			reqCtx := &model.RequestContext{}
			response := svc.Register(reqCtx, tt.userName, tt.password)

			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, tt.expectedOK, response.OK())
			if tt.expectedError != "" {
				assert.Contains(t, response.ErrorString(), tt.expectedError)
			}
			if tt.expectedOK {
				assert.NotNil(t, response.Data)
				assert.Empty(t, response.Data.Password, "Password should be removed from response")
			}

			mockRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name           string
		userName       string
		password       string
		setupMocks     func(*mocks.MockRepository, *mocks.MockUserRepo)
		expectedCode   int
		expectedOK     bool
		expectedError  string
		checkToken     bool
	}{
		{
			name:     "successful login",
			userName: "testuser",
			password: "password123",
			setupMocks: func(mockRepo *mocks.MockRepository, mockUserRepo *mocks.MockUserRepo) {
				hashedPassword, _ := hashPassword("password123")
				mockUserRepo.On("QueryOne", mock.Anything, &model.User{UserName: "testuser"}).
					Return(model.SuccessResponse(&model.User{
						ID:       1,
						UserName: "testuser",
						Password: hashedPassword,
					}, "User found"))
				mockRepo.On("User").Return(mockUserRepo)
			},
			expectedCode: model.CodeSuccess,
			expectedOK:   true,
			checkToken:   true,
		},
		{
			name:     "empty username",
			userName: "",
			password: "password123",
			setupMocks: func(mockRepo *mocks.MockRepository, mockUserRepo *mocks.MockUserRepo) {
				// No repo calls expected, service returns early
			},
			expectedCode:  model.CodeBadRequest,
			expectedOK:    false,
			expectedError: "Username and password are required",
		},
		{
			name:     "user not found",
			userName: "nonexistent",
			password: "password123",
			setupMocks: func(mockRepo *mocks.MockRepository, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.On("QueryOne", mock.Anything, &model.User{UserName: "nonexistent"}).
					Return(model.NotFound[*model.User]("User not found"))
				mockRepo.On("User").Return(mockUserRepo)
			},
			expectedCode:  model.CodeUnauthorized,
			expectedOK:    false,
			expectedError: "Invalid credentials",
		},
		{
			name:     "wrong password",
			userName: "testuser",
			password: "wrongpassword",
			setupMocks: func(mockRepo *mocks.MockRepository, mockUserRepo *mocks.MockUserRepo) {
				hashedPassword, _ := hashPassword("password123")
				mockUserRepo.On("QueryOne", mock.Anything, &model.User{UserName: "testuser"}).
					Return(model.SuccessResponse(&model.User{
						ID:       1,
						UserName: "testuser",
						Password: hashedPassword,
					}, "User found"))
				mockRepo.On("User").Return(mockUserRepo)
			},
			expectedCode:  model.CodeUnauthorized,
			expectedOK:    false,
			expectedError: "Invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockRepository)
			mockUserRepo := new(mocks.MockUserRepo)
			tt.setupMocks(mockRepo, mockUserRepo)

			svc := auth.NewTestAuthService(mockRepo, "test-secret-key-for-testing-only")

			reqCtx := &model.RequestContext{}
			response := svc.Login(reqCtx, tt.userName, tt.password)

			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, tt.expectedOK, response.OK())
			if tt.expectedError != "" {
				assert.Contains(t, response.ErrorString(), tt.expectedError)
			}
			if tt.checkToken {
				assert.NotEmpty(t, response.Data, "Token should be generated")
			}

			mockRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Authenticate(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		setupMocks    func(auth.AuthService)
		expectedCode  int
		expectedOK    bool
		expectedError string
	}{
		{
			name:  "missing token",
			token: "",
			setupMocks: func(svc auth.AuthService) {
				// No mocks needed
			},
			expectedCode:  model.CodeUnauthorized,
			expectedOK:    false,
			expectedError: "Token is required",
		},
		{
			name:  "invalid token",
			token: "invalid-token",
			setupMocks: func(svc auth.AuthService) {
				// No mocks needed, ParseToken will fail
			},
			expectedCode:  model.CodeUnauthorized,
			expectedOK:    false,
			expectedError: "Invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := auth.NewTestAuthService(nil, "test-secret-key-for-testing-only")
			tt.setupMocks(svc)

			reqCtx := &model.RequestContext{Token: tt.token}
			response := svc.Authenticate(reqCtx)

			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, tt.expectedOK, response.OK())
			if tt.expectedError != "" {
				assert.Contains(t, response.ErrorString(), tt.expectedError)
			}
		})
	}
}

func TestAuthService_GetMe(t *testing.T) {
	// Generate a valid token for testing
	jwtSecret := "test-secret-key-for-testing-only"
	claims := &auth.TestJWTClaims{
		SessionID: uuid.New().String(),
		UserID:    1,
		UserName:  "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, _ := token.SignedString([]byte(jwtSecret))

	tests := []struct {
		name          string
		token         string
		setupMocks    func(*mocks.MockRepository, *mocks.MockUserRepo)
		expectedCode  int
		expectedOK    bool
	}{
		{
			name:  "successful get me",
			token: validToken,
			setupMocks: func(mockRepo *mocks.MockRepository, mockUserRepo *mocks.MockUserRepo) {
				mockUserRepo.On("QueryOne", mock.Anything, &model.User{ID: uint(1)}).
					Return(model.SuccessResponse(&model.User{
						ID:       1,
						UserName: "testuser",
						Password: "hashed",
					}, "User found"))
				mockRepo.On("User").Return(mockUserRepo)
			},
			expectedCode: model.CodeSuccess,
			expectedOK:   true,
		},
		{
			name:  "missing token",
			token: "",
			setupMocks: func(mockRepo *mocks.MockRepository, mockUserRepo *mocks.MockUserRepo) {
				// No repo calls expected
			},
			expectedCode: model.CodeUnauthorized,
			expectedOK:   false,
		},
		{
			name:  "invalid token",
			token: "invalid-token",
			setupMocks: func(mockRepo *mocks.MockRepository, mockUserRepo *mocks.MockUserRepo) {
				// No repo calls expected, ParseToken will fail
			},
			expectedCode: model.CodeUnauthorized,
			expectedOK:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockRepository)
			mockUserRepo := new(mocks.MockUserRepo)
			tt.setupMocks(mockRepo, mockUserRepo)

			svc := auth.NewTestAuthService(mockRepo, jwtSecret)

			reqCtx := &model.RequestContext{Token: tt.token}
			response := svc.GetMe(reqCtx)

			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, tt.expectedOK, response.OK())
			if tt.expectedOK {
				assert.NotNil(t, response.Data)
				assert.Empty(t, response.Data.Password, "Password should be removed")
			}

			mockRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

// Helper functions
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func hashPasswordForTest(password string) string {
	hashed, _ := hashPassword(password)
	return hashed
}

