package socket

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	Event   string `json:"event"`
	Payload any    `json:"payload"`
}

func RunEngine(socket Socket) {
	socket.DispatchEvent("connected", nil)
	for {
		_, content, err := socket.GetConnect().ReadMessage()
		if err != nil {
			socket.Disconnect()
			return
		}

		data := Request{}
		err = json.Unmarshal(content, &data)

		if err == nil {
			socket.DispatchEvent(data.Event, data.Payload)
		} else {
			fmt.Print(err)
		}
	}
}
