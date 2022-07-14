package wss

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	jsoniter "github.com/json-iterator/go"
)

// webSocket请求ping 返回pong
func WsServer() func(*fiber.Ctx) error {
	return websocket.New(func(ws *websocket.Conn) {
		AuthStatus := -1

		var (
			mt  int
			msg []byte
			err error
		)

		go func() {
			for {
				mt, msg, err = ws.ReadMessage()
				if err != nil {
					ws.Close()
					break
				}

				SendData := Read(msg)
				AuthStatus = SendData.Code

				b, _ := jsoniter.Marshal(SendData)
				ws.WriteMessage(mt, b)
				time.Sleep(time.Second) // 一秒执行一次
			}
		}()

		for {
			SendData := Send()

			if AuthStatus > 0 {
				b, _ := jsoniter.Marshal(SendData)

				err := ws.WriteMessage(1, b)
				if err != nil {
					ws.Close()
					break
				}
			}

			time.Sleep(time.Second) // 一秒执行一次
		}
	})
}
