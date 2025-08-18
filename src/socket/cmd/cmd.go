package cmd

import (
	"fmt"
	"local/config"
	"local/event"
	"local/libs/socket"
	Router "local/router"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func Run() {
	config.LoadConfig()
	router := http.NewServeMux()

	socketServer := socket.NewServer(router)
	event.RegisterEvent(socketServer)
	Router.Register(router, socketServer)
	log.Println(fmt.Sprintf("Server is running on host %s and port %d", config.Config.Host, config.Config.HTTPPort))
	http.ListenAndServe(config.Config.Host+":"+strconv.Itoa(config.Config.HTTPPort), router)
}
