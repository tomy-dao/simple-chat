package socket

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type SocketMap map[string]Socket

type SocketStore interface {
	Add(socket Socket)
	Get(socketId uuid.UUID) (Socket, error)
	RemoveID(socketId uuid.UUID)
	Remove(socket Socket)
	GetAll() SocketMap
}

type socketStore struct {
	sockets SocketMap
	lock    sync.Mutex
}

func (s *socketStore) Add(socket Socket) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.sockets[socket.GetId()] = socket
}

func (s *socketStore) Get(socketId uuid.UUID) (Socket, error) {
	socket, ok := s.sockets[socketId.String()]
	if !ok {
		return nil, errors.New("not_found")
	}

	return socket, nil
}

func (s *socketStore) RemoveID(socketId uuid.UUID) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.sockets, socketId.String())
}

func (s *socketStore) Remove(socket Socket) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.sockets, socket.GetId())
}

func (s *socketStore) GetAll() SocketMap {
	return s.sockets
}

func NewSocketStore() SocketStore {
	sockets := make(SocketMap)
	return &socketStore{
		sockets: sockets,
	}
}
