package wss

import (
	"time"

	"github.com/EasyGolang/goTools/mRes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	jsoniter "github.com/json-iterator/go"
)

// webSocket请求ping 返回pong
func WsServer() func(*fiber.Ctx) error {
	return websocket.New(func(ws *websocket.Conn) {
		AuthResult := mRes.ResType{}

		var (
			msg []byte
			err error
		)

		go func() {
			for {
				_, msg, err = ws.ReadMessage()
				if err != nil {
					ws.Close()
					break
				}

				// 第一次需要验证
				if AuthResult.Code == 0 {
					AuthResult = Auth(msg)
				}

				time.Sleep(time.Second) // 一秒执行一次
			}
		}()

		for {

			if AuthResult.Code > 0 {
				AuthResult := Send()
				b, _ := jsoniter.Marshal(AuthResult)
				err := ws.WriteMessage(1, b)
				if err != nil {
					ws.Close()
					break
				}
			} else {
				b, _ := jsoniter.Marshal(AuthResult)
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
