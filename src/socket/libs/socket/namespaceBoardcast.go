package socket

import (
	"sync"

	"github.com/google/uuid"
)

var DefaultNamespace = "/"

type SocketSet map[string]struct{}

type NamespaceBroadcast interface {
	Of(path string) NamespaceBroadcast
	ToRoom(room string) NamespaceBroadcast
	ToRooms(room []string) NamespaceBroadcast
	Emit(event string, message any)
	RemoveRoom(room string)
	GetSocketSelected() SocketMap
	To(socketID uuid.UUID) NamespaceBroadcast
	WithoutRoom(room string) NamespaceBroadcast
	WithoutConn(conns string) NamespaceBroadcast
}

func NewNamespaceBroadcast(store NamespaceStore) NamespaceBroadcast {
	roomSet := make(RoomSet)
	socketSet := make(SocketSet)
	return &namespaceBroadcast{
		namespaceStore: store,
		namespace:      DefaultNamespace,
		roomSet:        roomSet,
		socketSet:      socketSet,
	}
}

type namespaceBroadcast struct {
	namespaceStore NamespaceStore
	namespace      string
	roomSet        RoomSet
	socketSet      SocketSet
	withoutRooms   []string
	withoutConns   []string
	lock           sync.RWMutex
}

func (n *namespaceBroadcast) RemoveRoom(room string) {
	skRoom := n.namespaceStore.Get(n.namespace).Get(room)
	if skRoom == nil {
		return
	}
	for _, sk := range skRoom.GetAll() {
		sk.LeaveRom(room)
	}
	n.namespaceStore.Get(n.namespace).Remove(room)
}

func (n *namespaceBroadcast) Of(path string) NamespaceBroadcast {
	socketSet := make(SocketSet)

	return &namespaceBroadcast{
		namespaceStore: n.namespaceStore,
		namespace:      path,
		roomSet:        n.roomSet,
		socketSet:      socketSet,
	}
}

func (n *namespaceBroadcast) ToRoom(room string) NamespaceBroadcast {
	roomSet := copyRoomSet(n.roomSet)
	roomSet[room] = empty

	return &namespaceBroadcast{
		namespaceStore: n.namespaceStore,
		namespace:      n.namespace,
		roomSet:        roomSet,
		socketSet:      n.socketSet,
	}
}

func (n *namespaceBroadcast) ToRooms(rooms []string) NamespaceBroadcast {
	roomSet := copyRoomSet(n.roomSet)

	for _, room := range rooms {
		roomSet[room] = empty
	}

	return &namespaceBroadcast{
		namespaceStore: n.namespaceStore,
		namespace:      n.namespace,
		roomSet:        roomSet,
		socketSet:      n.socketSet,
	}
}

func (n *namespaceBroadcast) Emit(event string, message any) {
	skMap := n.GetSocketSelected()

	withoutConnectIds := make(map[string]struct{})
	for _, conn := range n.withoutConns {
		if conn != "" {
			withoutConnectIds[conn] = empty
		}
	}

	for _, room := range n.withoutRooms {
		skStore := n.namespaceStore.
			Get(n.namespace).
			Get(room)
		if skStore == nil {
			continue
		}

		conns := skStore.GetAll()
		for _, conn := range conns {
			if conn != nil {
				withoutConnectIds[conn.GetId()] = empty
			}
		}
	}

	if skMap == nil {
		return
	}

	for _, sk := range skMap {
		if sk == nil {
			continue
		}
		if _, ok := withoutConnectIds[sk.GetId()]; ok {
			continue
		}
		sk.Emit(event, message)
	}
}

func (n *namespaceBroadcast) GetSocketSelected() SocketMap {
	rStore := n.namespaceStore.Get(n.namespace)
	if rStore == nil {
		return nil
	}

	skStore := NewSocketStore()

	// get socket set
	for skId, _ := range n.socketSet {
		id, err := uuid.Parse(skId)
		if err != nil {
			continue
		}

		sk := rStore.GetSocket(id)
		if sk == nil {
			continue
		}

		skStore.Add(sk)
	}

	// get socket in room set
	for room, _ := range n.roomSet {
		sockets := rStore.Get(room)
		if sockets == nil {
			continue
		}

		for _, sk := range sockets.GetAll() {
			if sk != nil {
				skStore.Add(sk)
			}
		}
	}

	// return all socket of namespace when haven't set room and specify socket
	if len(n.socketSet) == 0 && len(n.roomSet) == 0 {
		return rStore.Get(DefaultRoom).GetAll()
	}
	// return sockets when have set room and specify socket
	return skStore.GetAll()
}

func (n *namespaceBroadcast) WithoutRoom(room string) NamespaceBroadcast {
	n.lock.Lock()
	withoutRooms := append(n.withoutRooms, room)
	n.lock.Unlock()

	return &namespaceBroadcast{
		namespaceStore: n.namespaceStore,
		namespace:      n.namespace,
		roomSet:        n.roomSet,
		socketSet:      n.socketSet,
		withoutRooms:   withoutRooms,
	}
}

func (n *namespaceBroadcast) WithoutConn(conn string) NamespaceBroadcast {
	n.lock.Lock()
	withoutConns := append(n.withoutConns, conn)
	n.lock.Unlock()

	return &namespaceBroadcast{
		namespaceStore: n.namespaceStore,
		namespace:      n.namespace,
		roomSet:        n.roomSet,
		socketSet:      n.socketSet,
		withoutConns:   withoutConns,
	}
}

func (n *namespaceBroadcast) To(socketID uuid.UUID) NamespaceBroadcast {
	socketSet := copySocketSet(n.socketSet)
	socketSet[socketID.String()] = empty

	return &namespaceBroadcast{
		namespaceStore: n.namespaceStore,
		namespace:      n.namespace,
		roomSet:        n.roomSet,
		socketSet:      socketSet,
	}
}
