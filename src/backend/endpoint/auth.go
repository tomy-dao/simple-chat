package endpoint

import (
	"local/model"
	"local/service/auth"
	"local/service/initial"
	"local/util/logger"
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

func (e *AuthEndpoints) Authenticate(reqCtx *model.RequestContext) model.Response[uint] {
	logger.Info(reqCtx, "AuthEndpoints.Authenticate called")
	return e.authService.Authenticate(reqCtx)
}

func (e *AuthEndpoints) GetMe(reqCtx *model.RequestContext) model.Response[*model.User] {
	logger.Info(reqCtx, "AuthEndpoints.GetMe called")
	return e.authService.GetMe(reqCtx)
}

func (e *AuthEndpoints) Register(reqCtx *model.RequestContext, request interface{}) model.Response[*model.User] {
	req := request.(RegisterRequest)
	logger.Info(reqCtx, "AuthEndpoints.Register called", map[string]interface{}{"username": req.UserName})
	return e.authService.Register(reqCtx, req.UserName, req.Password)
}

func (e *AuthEndpoints) Login(reqCtx *model.RequestContext, request interface{}) model.Response[LoginResponse] {
	req := request.(LoginRequest)
	logger.Info(reqCtx, "AuthEndpoints.Login called", map[string]interface{}{"username": req.UserName})

	tokenResponse := e.authService.Login(reqCtx, req.UserName, req.Password)
	if !tokenResponse.OK() {
		return model.ErrorArray[LoginResponse](tokenResponse.Code, tokenResponse.Message, tokenResponse.Errors)
	}

	response := LoginResponse{Token: tokenResponse.Data}
	return model.SuccessResponse(response, "Login successful")
}

func (e *AuthEndpoints) Logout(reqCtx *model.RequestContext, request interface{}) model.Response[string] {
	req := request.(LogoutRequest)
	logger.Info(reqCtx, "AuthEndpoints.Logout called")
	return e.authService.Logout(reqCtx, req.Token)
}

func (e *AuthEndpoints) GetUsers(reqCtx *model.RequestContext) model.Response[[]*model.User] {
	logger.Info(reqCtx, "AuthEndpoints.GetUsers called")
	return e.authService.GetUsers(reqCtx)
}

func NewAuthEndpoints(params *initial.Service) *AuthEndpoints {
	return &AuthEndpoints{
		authService: params.AuthSvc,
	}
}
