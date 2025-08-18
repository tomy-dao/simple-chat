package router

import (
	"local/handler"
	SK "local/libs/socket"
	"net/http"
)

func Register(r *http.ServeMux, socketServer SK.Server) {
	handler := handler.NewHandler(socketServer)

	r.HandleFunc("/broadcast", handler.Broadcast)
	
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"health_check":"OK"}`))
	})
}
