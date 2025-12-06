package auth

import (
	"fmt"
	"local/config"
	"local/infra/repo"
	"local/model"
	"local/service/common"
	"local/util/logger"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims represents the JWT token claims structure
type JWTClaims struct {
	SessionID string `json:"session_id"`
	UserID    uint   `json:"user_id"`
	UserName  string `json:"user_name"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Authenticate(reqCtx *model.RequestContext) model.Response[uint]
	ParseToken(tokenStr string) model.Response[*JWTClaims]
	CheckToken(reqCtx *model.RequestContext, token string) model.Response[bool]
	GetMe(reqCtx *model.RequestContext) model.Response[*model.User]
	Register(reqCtx *model.RequestContext, userName, password string) model.Response[*model.User]
	Login(reqCtx *model.RequestContext, userName, password string) model.Response[string]
	Logout(reqCtx *model.RequestContext, token string) model.Response[string]
	GetUsers(reqCtx *model.RequestContext) model.Response[[]*model.User]
}

type authService struct {
	repo      repo.RepositoryInterface
	jwtSecret string
}

func (svc *authService) Authenticate(reqCtx *model.RequestContext) model.Response[uint] {
	logger.Info(reqCtx, "Authenticate called")
	if reqCtx.Token == "" {
		return model.Unauthorized[uint]("Token is required")
	}

	tokenResponse := svc.ParseToken(reqCtx.Token)
	if !tokenResponse.OK() {
		return model.Unauthorized[uint](tokenResponse.ErrorString())
	}

	claims := tokenResponse.Data
	return model.SuccessResponse(claims.UserID, "Authentication successful")
}

func (svc *authService) ParseToken(tokenStr string) model.Response[*JWTClaims] {
	logger.Info(nil, "ParseToken called", map[string]interface{}{"token_length": len(tokenStr)})
	// Parse token
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra thuật toán
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(svc.jwtSecret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return model.Unauthorized[*JWTClaims]("Invalid token")
	}

	// Lấy claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return model.SuccessResponse(claims, "Token parsed successfully")
	}

	return model.Unauthorized[*JWTClaims]("Invalid token claims")
}

func (svc *authService) CheckToken(reqCtx *model.RequestContext, token string) model.Response[bool] {
	logger.Info(reqCtx, "CheckToken called", map[string]interface{}{"token_length": len(token)})
	if token == "" {
		return model.BadRequest[bool]("Token is required")
	}

	tokenResponse := svc.ParseToken(token)
	if !tokenResponse.OK() {
		return model.Unauthorized[bool](tokenResponse.ErrorString())
	}
	return model.SuccessResponse(true, "Token is valid")
}

func (svc *authService) GetMe(reqCtx *model.RequestContext) model.Response[*model.User] {
	logger.Info(reqCtx, "GetMe called")
	if reqCtx.Token == "" {
		return model.Unauthorized[*model.User]("Token is required")
	}
	
	tokenResponse := svc.ParseToken(reqCtx.Token)
	if !tokenResponse.OK() {
		return model.Unauthorized[*model.User](tokenResponse.ErrorString())
	}

	claims := tokenResponse.Data
	response := svc.repo.User().QueryOne(reqCtx, &model.User{ID: claims.UserID})
	if !response.OK() {
		return response
	}

	// Remove password from response
	response.Data.Password = ""
	return response
}

func (svc *authService) Register(reqCtx *model.RequestContext, userName, password string) model.Response[*model.User] {
	logger.Info(reqCtx, "Register called", map[string]interface{}{"username": userName})
	if userName == "" || password == "" {
		return model.BadRequest[*model.User]("Username and password are required")
	}
	// Check if user already exists
	existingUserResponse := svc.repo.User().QueryOne(reqCtx, &model.User{UserName: userName})
	if existingUserResponse.OK() {
		return model.Conflict[*model.User]("User already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.InternalError[*model.User]("Failed to hash password")
	}

	// Create user
	user := &model.User{
		UserName: userName,
		Password: string(hashedPassword),
	}

	response := svc.repo.User().Create(reqCtx, user)
	if !response.OK() {
		return response
	}

	// Remove password from response
	response.Data.Password = ""
	return response
}

func (svc *authService) Login(reqCtx *model.RequestContext, userName, password string) model.Response[string] {
	logger.Info(reqCtx, "Login called", map[string]interface{}{"username": userName})
	if userName == "" || password == "" {
		return model.BadRequest[string]("Username and password are required")
	}

	// Find user
	response := svc.repo.User().QueryOne(reqCtx, &model.User{UserName: userName})
	if !response.OK() {
		return model.Unauthorized[string]("Invalid credentials")
	}

	user := response.Data

	// Check password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return model.Unauthorized[string]("Invalid credentials")
	}

	// Generate JWT token
	claims := &JWTClaims{
		SessionID: uuid.New().String(),
		UserID:    user.ID,
		UserName:  user.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(svc.jwtSecret))
	if err != nil {
		return model.InternalError[string]("Failed to generate token")
	}

	return model.SuccessResponse(tokenString, "Login successful")
}

func (svc *authService) Logout(reqCtx *model.RequestContext, token string) model.Response[string] {
	logger.Info(reqCtx, "Logout called", map[string]interface{}{"token_length": len(token)})
	// In a real implementation, you might want to blacklist the token
	// For now, we'll just validate that the token is valid
	if token == "" {
		return model.BadRequest[string]("Token is required")
	}

	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(svc.jwtSecret), nil
	})

	if err != nil {
		return model.Unauthorized[string]("Invalid token")
	}

	// Token is valid, logout successful
	// In a production system, you'd typically add the token to a blacklist
	return model.SuccessResponse("", "Logout successful")
}

func (svc *authService) GetUsers(reqCtx *model.RequestContext) model.Response[[]*model.User] {
	logger.Info(reqCtx, "GetUsers called")
	response := svc.repo.User().QueryMany(reqCtx, &model.User{})
	return response
}

func NewAuthService(params *common.Params) AuthService {
	return &authService{
		repo:      params.Repo,
		jwtSecret: config.Config.JwtSecret,
	}
}

