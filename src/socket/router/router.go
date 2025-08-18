package router

import (
	"local/handler"
	SK "local/libs/socket"
	"net/http"
)

func Register(r *http.ServeMux, socketServer SK.Server) {
	handler := handler.NewHandler(socketServer)
	r.HandleFunc("/broadcast", handler.Broadcast)
}
