package auth

import (
	"context"
	"errors"
	"fmt"
	"local/config"
	"local/model"
	"local/repository"
	"local/service/common"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	Authenticate(ctx context.Context) (uint, error)
	ParseToken(tokenStr string) (map[string]interface{}, error)
	CheckToken(ctx context.Context, token string) (bool, error)
	GetMe(ctx context.Context) (*model.User, error)
	Register(ctx context.Context, userName, password string) (*model.User, error)
	Login(ctx context.Context, userName, password string) (string, error)
	Logout(ctx context.Context, token string) error
	GetUsers(ctx context.Context) ([]*model.User, error)
}

type authService struct {
	repo      repository.RepositoryInterface
	jwtSecret string
}

func (svc *authService) Authenticate(ctx context.Context) (uint, error) {
	token := ctx.Value("token").(string)
	if token == "" {
		return 0, errors.New("token is required")
	}

	valid, err := svc.ParseToken(token)
	if err != nil {
		return 0, err
	}

	userID := valid["user_id"].(float64)

	return uint(userID), nil
}

func (svc *authService) ParseToken(tokenStr string) (map[string]interface{}, error) {
	// Parse token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra thuật toán
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(svc.jwtSecret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Lấy claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Kiểm tra token có hết hạn không
		if exp, ok := claims["exp"]; ok {
			if expTime, ok := exp.(float64); ok {
				if time.Now().Unix() > int64(expTime) {
					return nil, errors.New("token expired")
				}
			}
		}
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

func (svc *authService) CheckToken(ctx context.Context, token string) (bool, error) {
	if token == "" {
		return false, errors.New("token is required")
	}

	_, err := svc.ParseToken(token)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (svc *authService) GetMe(ctx context.Context) (*model.User, error) {
	token := ctx.Value("token").(string)
	if token == "" {
		return nil, errors.New("token is required")
	}
	
	claims, err := svc.ParseToken(token)
	if err != nil {
		return nil, err
	}

	userID	:= claims["user_id"].(float64)

	user := svc.repo.User().QueryOne(ctx, &model.User{ID: uint(userID)})
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (svc *authService) Register(ctx context.Context, userName, password string) (*model.User, error) {
	if userName == "" || password == "" {
		return nil, errors.New("username and password are required")
	}
	// Check if user already exists
	existingUser := svc.repo.User().QueryOne(ctx, &model.User{UserName: userName})
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &model.User{
		UserName: userName,
		Password: string(hashedPassword),
	}

	createdUser := svc.repo.User().Create(ctx, user)
	if createdUser == nil {
		return nil, errors.New("failed to create user")
	}

	// Remove password from response
	createdUser.Password = ""
	return createdUser, nil
}

func (svc *authService) Login(ctx context.Context, userName, password string) (string, error) {
	if userName == "" || password == "" {
		return "", errors.New("username and password are required")
	}

	// Find user
	user := svc.repo.User().QueryOne(ctx, &model.User{UserName: userName})
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	// Check password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session_id": uuid.New().String(),
		"user_id":    user.ID,
		"user_name":  user.UserName,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":        time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(svc.jwtSecret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

func (svc *authService) Logout(ctx context.Context, token string) error {
	// In a real implementation, you might want to blacklist the token
	// For now, we'll just validate that the token is valid
	if token == "" {
		return errors.New("token is required")
	}

	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(svc.jwtSecret), nil
	})

	if err != nil {
		return errors.New("invalid token")
	}

	// Token is valid, logout successful
	// In a production system, you'd typically add the token to a blacklist
	return nil
}

func (svc *authService) GetUsers(ctx context.Context) ([]*model.User, error) {
	users := svc.repo.User().QueryMany(ctx, &model.User{})
	return users, nil
}

func NewAuthService(params *common.Params) AuthService {
	return &authService{
		repo:      params.Repo,
		jwtSecret: config.Config.JwtSecret,
	}
}
