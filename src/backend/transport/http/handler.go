package httpTransoprt

import (
	"encoding/json"
	"fmt"
	"local/endpoint"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	endpoints *endpoint.Endpoints
}

func getToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	fmt.Println("token", token)
	if token == "" {
		return ""
	}
	return token[7:]
}

func (h *handler) GetMe() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)
		req := endpoint.GetMeRequest{
			Token: token,
		}

		user, err := h.endpoints.Auth.GetMe(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

// Auth handlers
func (h *handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req endpoint.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		fmt.Println("Register request: ", req)
		user, err := h.endpoints.Auth.Register(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func (h *handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req endpoint.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := h.endpoints.Auth.Login(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(response)
	}
}

func (h *handler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)
		response, err := h.endpoints.Auth.Logout(r.Context(), endpoint.LogoutRequest{Token: token})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(response)
	}
}

type CreateConversationRequest struct {
	UserID uint `json:"user_id"`
}

// Conversation handlers
func (h *handler) CreateConversation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateConversationRequest
		userID, err := h.endpoints.Auth.Authenticate(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		conversation, err := h.endpoints.Conversation.CreateConversation(r.Context(), []uint{req.UserID, userID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(conversation)
	}
}

func (h *handler) GetConversationByUserID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		myID, err := h.endpoints.Auth.Authenticate(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		conversation, err := h.endpoints.Conversation.GetConversationByUserIDs(r.Context(), []uint{myID, uint(userID)})
		fmt.Println("Conversation: ", conversation)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(conversation)
	}
}

func (h *handler) GetConversations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := h.endpoints.Auth.Authenticate(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		conversations, err := h.endpoints.Conversation.GetUserConversations(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(conversations)
	}
}
// Message handlers
func (h *handler) CreateMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req endpoint.CreateMessageRequest
		userID, err := h.endpoints.Auth.Authenticate(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		conversationID, err := strconv.ParseUint(chi.URLParam(r, "conversationID"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
			return
		}
		req.ConversationID = uint(conversationID)
		req.SenderID = userID
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		message, err := h.endpoints.Message.CreateMessage(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(message)
	}
}

func (h *handler) GetMessagesByConversationID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conversationID, err := strconv.ParseUint(chi.URLParam(r, "conversationID"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
			return
		}
		messages, err := h.endpoints.Message.GetMessagesByConversationID(r.Context(), uint(conversationID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(messages)
	}
}

// User handlers
func (h *handler) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := h.endpoints.Auth.Authenticate(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		users, err := h.endpoints.Auth.GetUsers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}
