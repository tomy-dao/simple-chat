package endpoint

import (
	"context"
	"local/model"
	"local/service/auth"
	"local/service/initial"
)

type AuthEndpoints struct {
	authService auth.AuthService
}

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LogoutRequest struct {
	Token string `json:"token"`
}

type GetMeRequest struct {
	Token string `json:"token"`
}

func (e *AuthEndpoints) Authenticate(ctx context.Context) (uint, error) {
	userID, err := e.authService.Authenticate(ctx)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (e *AuthEndpoints) GetMe(ctx context.Context, request interface{}) (Response[*model.User], error) {
	req := request.(GetMeRequest)
	
	user, err := e.authService.GetMe(ctx, req.Token)
	if err != nil {
		return Response[*model.User]{Data: nil, Error: err.Error()}, nil
	}
	
	return Response[*model.User]{Data: &user, Error: ""}, nil
}

func (e *AuthEndpoints) Register(ctx context.Context, request interface{}) (Response[*model.User], error) {
	req := request.(RegisterRequest)

	user, err := e.authService.Register(ctx, req.UserName, req.Password)
	if err != nil {
		return Response[*model.User]{Data: nil, Error: err.Error()}, nil
	}

	return Response[*model.User]{Data: &user, Error: ""}, nil
}

func (e *AuthEndpoints) Login(ctx context.Context, request interface{}) (Response[LoginResponse], error) {
	req := request.(LoginRequest)

	token, err := e.authService.Login(ctx, req.UserName, req.Password)
	if err != nil {
		return Response[LoginResponse]{Data: nil, Error: err.Error()}, nil
	}

	response := LoginResponse{Token: token}
	return Response[LoginResponse]{Data: &response, Error: ""}, nil
}

func (e *AuthEndpoints) Logout(ctx context.Context, request interface{}) (Response[string], error) {
	req := request.(LogoutRequest)

	err := e.authService.Logout(ctx, req.Token)
	if err != nil {
		return Response[string]{Data: nil, Error: err.Error()}, nil
	}

	return Response[string]{Data: nil, Error: ""}, nil
}

func (e *AuthEndpoints) GetUsers(ctx context.Context) (Response[[]*model.User], error) {
	users, err := e.authService.GetUsers(ctx)
	if err != nil {
		return Response[[]*model.User]{Data: nil, Error: err.Error()}, nil
	}

	return Response[[]*model.User]{Data: &users, Error: ""}, nil
}

func NewAuthEndpoints(params *initial.Service) *AuthEndpoints {
	return &AuthEndpoints{
		authService: params.AuthSvc,
	}
}
