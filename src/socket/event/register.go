package event

import (
	"fmt"
	"local/client"
	SK "local/libs/socket"
)

type authenticate struct {
	Token string `json:"token"`
}

func RegisterEvent(socketServer SK.Server) {
	socketServer.Register(ChatPath, func(socket SK.Socket) {
		socket.On("connected", func(_ any) {
			socket.Emit("send_connect_id", socket.GetId())
		})
		
		socket.On("authenticate", func(payload any) {
			data := authenticate{}
			err := mapData(payload, &data)
			if err != nil {
				socket.Emit("authenticate_fail", nil)
				return
			}
			me, err := client.GetMe(data.Token)
			if err != nil {
				socket.Emit("authenticate_fail", nil)
				return
			}
			claims, err := decodeJWT(data.Token)
			if err != nil {
				socket.Emit("authenticate_fail", nil)
				return
			}

			sessionId := claims["session_id"].(string)
			userId, ok := claims["user_id"].(float64)
			if ok {
				socket.Join(sessionId)
				socket.Join(fmt.Sprintf("%d", int(userId)))
			}

			socket.Emit("authenticate_success", me)
		})
	})
}
