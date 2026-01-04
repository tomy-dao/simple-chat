package httpTransport

import (
	"local/endpoint"
	"local/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type handler struct {
	endpoints *endpoint.Endpoints
}

// GetMe godoc
// @Summary Get current authenticated user information
// @Description Returns the current user's information based on the JWT token
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} model.Response[model.User]
// @Failure 400 {object} model.Response[any] "Bad Request - Invalid input"
// @Failure 401 {object} model.Response[any] "Unauthorized - Invalid or missing token"
// @Failure 500 {object} model.Response[any] "Internal Server Error"
// @Router /me [get]
func (h *handler) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCtx := model.NewRequestContext(c.Request.Context())
		response := h.endpoints.Auth.GetMe(reqCtx)
		c.JSON(response.Code, response)
	}
}

// Register godoc
// @Summary Register a new user account
// @Description Creates a new user account with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body endpoint.RegisterRequest true "User registration data"
// @Success 200 {object} model.Response[model.User]
// @Failure 400 {object} model.Response[any] "Bad Request - Invalid input"
// @Failure 409 {object} model.Response[any] "Conflict - User already exists"
// @Failure 500 {object} model.Response[any] "Internal Server Error"
// @Router /register [post]
func (h *handler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req endpoint.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response := model.ValidationError[*model.User]("Invalid request body")
			c.JSON(response.Code, response)
			return
		}
		reqCtx := model.NewRequestContext(c.Request.Context())
		response := h.endpoints.Auth.Register(reqCtx, req)
		c.JSON(response.Code, response)
	}
}

// Login godoc
// @Summary Authenticate user and get JWT token
// @Description Validates user credentials and returns a JWT token for authentication
// @Tags auth
// @Accept json
// @Produce json
// @Param request body endpoint.LoginRequest true "User login credentials"
// @Success 200 {object} model.Response[endpoint.LoginResponse]
// @Failure 400 {object} model.Response[any] "Bad Request - Invalid input"
// @Failure 401 {object} model.Response[any] "Unauthorized - Invalid credentials"
// @Failure 500 {object} model.Response[any] "Internal Server Error"
// @Router /login [post]
func (h *handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req endpoint.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response := model.ValidationError[endpoint.LoginResponse]("Invalid request body")
			c.JSON(response.Code, response)
			return
		}

		reqCtx := model.NewRequestContext(c.Request.Context())
		response := h.endpoints.Auth.Login(reqCtx, req)
		c.JSON(response.Code, response)
	}
}

// Logout godoc
// @Summary Logout user and invalidate token
// @Description Invalidates the current user's JWT token
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.Response[string]
// @Failure 400 {object} model.Response[any] "Bad Request - Invalid input"
// @Failure 401 {object} model.Response[any] "Unauthorized - Invalid or missing token"
// @Failure 500 {object} model.Response[any] "Internal Server Error"
// @Router /logout [post]
func (h *handler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		reqCtx := model.NewRequestContext(c.Request.Context())
		response := h.endpoints.Auth.Logout(reqCtx, endpoint.LogoutRequest{Token: token})
		c.JSON(response.Code, response)
	}
}

type CreateConversationRequest struct {
	UserID uint `json:"user_id"`
}

// CreateConversation godoc
// @Summary Create a new private conversation
// @Description Creates a new private conversation between the authenticated user and another user
// @Tags conversations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateConversationRequest true "Conversation creation data"
// @Success 200 {object} model.Response[model.Conversation]
// @Failure 400 {object} model.Response[any] "Bad Request - Invalid input"
// @Failure 401 {object} model.Response[any] "Unauthorized - Invalid or missing token"
// @Failure 500 {object} model.Response[any] "Internal Server Error"
// @Router /conversations [post]
func (h *handler) CreateConversation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateConversationRequest
		reqCtx := model.NewRequestContext(c.Request.Context())
		if reqCtx.UserID == 0 {
			response := model.Unauthorized[*model.Conversation]("Unauthorized")
			c.JSON(response.Code, response)
			return
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response := model.ValidationError[*model.Conversation]("Invalid request body")
			c.JSON(response.Code, response)
			return
		}

		response := h.endpoints.Conversation.CreateConversation(reqCtx, []uint{req.UserID, reqCtx.UserID})
		c.JSON(response.Code, response)
	}
}

// GetConversationByUserID godoc
// @Summary Get conversation between current user and another user
// @Description Retrieves the private conversation between the authenticated user and the specified user
// @Tags conversations
// @Security BearerAuth
// @Produce json
// @Param userID path int true "Target User ID"
// @Success 200 {object} model.Response[model.Conversation]
// @Failure 400 {object} model.Response[any] "Bad Request - Invalid input"
// @Failure 401 {object} model.Response[any] "Unauthorized - Invalid or missing token"
// @Failure 404 {object} model.Response[any] "Not Found - Conversation not found"
// @Failure 500 {object} model.Response[any] "Internal Server Error"
// @Router /conversations/user/{userID} [get]
func (h *handler) GetConversationByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCtx := model.NewRequestContext(c.Request.Context())
		if reqCtx.UserID == 0 {
			response := model.Unauthorized[*model.Conversation]("Unauthorized")
			c.JSON(response.Code, response)
			return
		}
		userIDStr := c.Param("userID")
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			response := model.ValidationError[*model.Conversation]("Invalid user ID")
			c.JSON(response.Code, response)
			return
		}

		response := h.endpoints.Conversation.GetConversationByUserIDs(reqCtx, []uint{reqCtx.UserID, uint(userID)})
		c.JSON(response.Code, response)
	}
}

// GetConversations godoc
// @Summary Get all conversations for the authenticated user
// @Description Returns a list of all conversations where the authenticated user is a participant
// @Tags conversations
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.Response[[]model.Conversation]
// @Failure 400 {object} model.Response[any] "Bad Request - Invalid input"
// @Failure 401 {object} model.Response[any] "Unauthorized - Invalid or missing token"
// @Failure 500 {object} model.Response[any] "Internal Server Error"
// @Router /conversations [get]
func (h *handler) GetConversations() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCtx := model.NewRequestContext(c.Request.Context())
		if reqCtx.UserID == 0 {
			response := model.Unauthorized[[]*model.Conversation]("Unauthorized")
			c.JSON(response.Code, response)
			return
		}
		response := h.endpoints.Conversation.GetUserConversations(reqCtx, reqCtx.UserID)
		c.JSON(response.Code, response)
	}
}

// CreateMessage godoc
// @Summary Create a new message in a conversation
// @Description Creates a new message in the specified conversation
// @Tags messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param conversationID path int true "Conversation ID"
// @Param request body endpoint.CreateMessageRequest true "Message data"
// @Success 200 {object} model.Response[model.Message]
// @Failure 400 {object} model.Response[any] "Bad Request - Invalid input"
// @Failure 401 {object} model.Response[any] "Unauthorized - Invalid or missing token"
// @Failure 404 {object} model.Response[any] "Not Found - Conversation not found"
// @Failure 500 {object} model.Response[any] "Internal Server Error"
// @Router /conversations/{conversationID}/messages [post]
func (h *handler) CreateMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req endpoint.CreateMessageRequest
		reqCtx := model.NewRequestContext(c.Request.Context())
		if reqCtx.UserID == 0 {
			response := model.Unauthorized[*model.Message]("Unauthorized")
			c.JSON(response.Code, response)
			return
		}
		conversationIDStr := c.Param("conversationID")
		conversationID, err := strconv.ParseUint(conversationIDStr, 10, 64)
		if err != nil {
			response := model.ValidationError[*model.Message]("Invalid conversation ID")
			c.JSON(response.Code, response)
			return
		}
		req.ConversationID = uint(conversationID)
		req.SenderID = reqCtx.UserID
		if err := c.ShouldBindJSON(&req); err != nil {
			response := model.ValidationError[*model.Message]("Invalid request body")
			c.JSON(response.Code, response)
			return
		}

		response := h.endpoints.Message.CreateMessage(reqCtx, req)
		c.JSON(response.Code, response)
	}
}

// GetMessagesByConversationID godoc
// @Summary Get all messages in a conversation
// @Description Retrieves all messages from the specified conversation, ordered by creation time (newest first)
// @Tags messages
// @Security BearerAuth
// @Produce json
// @Param conversationID path int true "Conversation ID"
// @Success 200 {object} model.Response[[]model.Message]
// @Failure 400 {object} model.Response[any] "Bad Request - Invalid input"
// @Failure 401 {object} model.Response[any] "Unauthorized - Invalid or missing token"
// @Failure 500 {object} model.Response[any] "Internal Server Error"
// @Router /conversations/{conversationID}/messages [get]
func (h *handler) GetMessagesByConversationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		conversationIDStr := c.Param("conversationID")
		conversationID, err := strconv.ParseUint(conversationIDStr, 10, 64)
		if err != nil {
			response := model.ValidationError[[]*model.Message]("Invalid conversation ID")
			c.JSON(response.Code, response)
			return
		}
		reqCtx := model.NewRequestContext(c.Request.Context())
		response := h.endpoints.Message.GetMessagesByConversationID(reqCtx, uint(conversationID))
		c.JSON(response.Code, response)
	}
}

// GetUsers godoc
// @Summary Get all registered users
// @Description Returns a list of all registered users in the system
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.Response[[]model.User]
// @Failure 400 {object} model.Response[any] "Bad Request - Invalid input"
// @Failure 401 {object} model.Response[any] "Unauthorized - Invalid or missing token"
// @Failure 500 {object} model.Response[any] "Internal Server Error"
// @Router /users [get]
func (h *handler) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCtx := model.NewRequestContext(c.Request.Context())
		response := h.endpoints.Auth.GetUsers(reqCtx)
		c.JSON(response.Code, response)
	}
}
