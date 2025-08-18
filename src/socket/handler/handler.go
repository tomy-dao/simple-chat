package handler

import (
	"encoding/json"
	"fmt"
	"local/event"
	"local/libs/socket"
	"log"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

type RouterHandler interface {
	Broadcast(w http.ResponseWriter, r *http.Request)
}

type handle struct {
	socketServer socket.Server
}

type responseStatus struct {
	Ok bool
}

type responseError struct {
	Error string
}

type RequestBroadcast struct {
	UserIds []int `json:"user_ids" validate:"required"`
	SessionId string `json:"session_id" validate:"required"`
	Event   string `json:"event" validate:"required"`
	Payload any    `json:"payload"`
}

func (h *handle) responseJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	jsonStr, err := json.Marshal(data)
	if err != nil {
		log.Default().Print("pars json error ", err)
		fmt.Fprintf(w, `{"err": "%v"}`, err)
		return
	}
	w.Write(jsonStr)
}

func (h *handle) Broadcast(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	var res RequestBroadcast
	log.Default().Print("Broadcast")
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		fmt.Fprintf(w, `{"err": "%v"}`, err)
		w.WriteHeader(http.StatusBadRequest)
		log.Default().Print("decode error ", err)
		return
	}
	if err := validate.Struct(res); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.responseJSON(w, &responseError{Error: err.Error()})
		log.Default().Print("validate error ", err)
		return
	}

	rooms := make([]string, len(res.UserIds))
	for i, userId := range res.UserIds {
		rooms[i] = fmt.Sprintf("%d", userId)
	}

	h.socketServer.
		Broadcast().
		Of(event.ChatPath).
		ToRooms(rooms).
		WithoutConn(res.SessionId).
		Emit(res.Event, res.Payload)

	log.Default().Print("Broadcast success")
	h.responseJSON(w, responseStatus{
		Ok: true,
	})
}

func NewHandler(socketServer socket.Server) RouterHandler {
	return &handle{
		socketServer: socketServer,
	}
}
