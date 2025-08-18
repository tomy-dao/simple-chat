package socket

import (
	"fmt"
	"net/http"
)

type namespaceOnConnect func(socket Socket)

type Server interface {
	Register(path string, fn namespaceOnConnect)
	Broadcast() NamespaceBroadcast
}

var RoomDefault = ""

type server struct {
	router         *http.ServeMux
	namespaceStore NamespaceStore
}

func (s *server) Register(path string, fn namespaceOnConnect) {
	s.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		sk, err := NewSocket(w, r, s.namespaceStore, path)
		defer sk.Disconnect()

		if err != nil {
			fmt.Println(err)
			return
		}
		fn(sk)

		RunEngine(sk)
	})
}

func (s *server) Broadcast() NamespaceBroadcast {
	return NewNamespaceBroadcast(s.namespaceStore)
}

func NewServer(r *http.ServeMux) Server {
	nspStore := NewNamespaceStore()
	return &server{
		router:         r,
		namespaceStore: nspStore,
	}
}
