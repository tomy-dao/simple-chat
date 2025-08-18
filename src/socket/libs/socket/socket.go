package socket

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
)

type EventHandler func(payload any)

type RoomSet map[string]struct{}

var empty = struct{}{}

type StrictEventEmitter interface {
	On(event string, fn EventHandler)
	Emit(event string, message any) error
}

type Socket interface {
	StrictEventEmitter
	Join(rom string)
	Disconnect()
	GetConnect() Connect
	GetId() string
	DispatchEvent(event string, message any)
	LeaveRom(room string)
	GetRooms() []string
	AttachValue(key string, value any)
	GetValueAttach(key string) any
}

type socketHandler struct {
	connect        Connect
	namespaceStore NamespaceStore
	eventMap       map[string]EventHandler
	lock           sync.Mutex
	namespace      string
	inRooms        RoomSet
	ctx            context.Context
}

func (s *socketHandler) Join(room string) {
	s.inRooms[room] = empty
	s.namespaceStore.Add(s.namespace, room, s)
}

func (s *socketHandler) LeaveRom(room string) {
	s.namespaceStore.Get(s.namespace).LeaveRoom(s, room)
	delete(s.inRooms, room)
}

func (s *socketHandler) On(event string, fn EventHandler) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.eventMap[event] = fn
}

func (s *socketHandler) Emit(event string, message any) error {
	res := make(map[string]any)
	res["event"] = event
	res["payload"] = message
	content, err := json.Marshal(res)
	if err != nil {
		return err
	}

	s.connect.GetConnect().WriteMessage(1, content)
	return nil
}

func (s *socketHandler) Disconnect() {
	for _, room := range s.GetRooms() {
		s.LeaveRom(room)
	}
	s.DispatchEvent("disconnected", nil)
	s.connect.Close()
}

func (s *socketHandler) GetConnect() Connect {
	return s.connect
}

func (s *socketHandler) GetId() string {
	return s.connect.GetId()
}

func (s *socketHandler) DispatchEvent(event string, message any) {
	fn, ok := s.eventMap[event]
	if ok {
		fn(message)
	}
}

func (s *socketHandler) GetRooms() []string {
	rooms := make([]string, len(s.inRooms))
	for room, _ := range s.inRooms {
		rooms = append(rooms, room)
	}

	return rooms
}

func (s *socketHandler) AttachValue(key string, value any) {
	s.ctx = context.WithValue(s.ctx, key, value)
}

func (s *socketHandler) GetValueAttach(key string) any {
	return s.ctx.Value(key)
}

func NewSocket(w http.ResponseWriter, r *http.Request, store NamespaceStore, nspName string) (Socket, error) {
	connect, err := NewSocketConnect(w, r)
	if err != nil {
		return nil, err
	}
	eventMap := make(map[string]EventHandler)
	inRooms := make(map[string]struct{})
	ctx := context.Background()

	s := &socketHandler{
		namespaceStore: store,
		connect:        connect,
		namespace:      nspName,
		eventMap:       eventMap,
		inRooms:        inRooms,
		ctx:            ctx,
	}

	s.Join(RoomDefault)

	return s, nil
}
