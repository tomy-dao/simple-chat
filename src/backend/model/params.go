package model

import (
	"context"
)

type InitParams struct {
	ServiceName string
	Ctx         context.Context
}

type BroadcastMessage struct {
	UserIds []int `json:"user_ids"`
	SessionId string `json:"session_id"`
	Event   string `json:"event"`
	Payload any    `json:"payload"`
}
