package cmd

import (
	"local/event"
	"local/libs/socket"
	Router "local/router"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run() {
	router := http.NewServeMux()

	router.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"health_check":"OK"}`))
	})

	socketServer := socket.NewServer(router)
	event.RegisterEvent(socketServer)
	Router.Register(router, socketServer)

	viper.AutomaticEnv()
	port := "8080"
	host := "0.0.0.0"
	// port := viper.GetString("PORT")
	// host := viper.GetString("HOST")
	log.Println("Server running", port)
	http.ListenAndServe(host+":"+port, router)
}
