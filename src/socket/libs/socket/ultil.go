package socket

func copyRoomSet(m RoomSet) RoomSet {
	mc := make(RoomSet)

	for key, value := range m {
		mc[key] = value
	}

	return mc
}

func copySocketSet(m SocketSet) SocketSet {
	mc := make(SocketSet)

	for key, value := range m {
		mc[key] = value
	}

	return mc
}
