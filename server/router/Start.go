package router

import (
	"os"
	"time"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/middle"
	"CoinAI.net/server/router/api"
	"CoinAI.net/server/router/api/sys"
	"CoinAI.net/server/router/wss"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Start() {
	// 加载日志文件
	fileName := config.Dir.Log + "/HTTP-T" + time.Now().Format("06年1月02日15时") + ".log"
	logFile, _ := os.Create(fileName)
	/*
		加载模板
		https://www.gouguoyin.cn/posts/10103.html
	*/

	// 创建服务
	app := fiber.New(fiber.Config{
		ServerHeader: "CoinAI.net",
	})

	app.Use(
		limiter.New(limiter.Config{
			Max:        100,
			Expiration: 1 * time.Second,
		}), // 限流
		logger.New(logger.Config{
			Format:     "[${time}] [${ip}:${port}] ${status} - ${method} ${latency} ${path} \n",
			TimeFormat: "2006-01-02 - 15:04:05",
			Output:     logFile,
		}), // 日志
		cors.New(),     // 允许跨域
		compress.New(), // 压缩
		middle.Public,  // 授权验证
	)

	// CoinAI
	r_api := app.Group("/CoinAI")
	r_api.Get("/config", api.GetConfig)
	r_api.Get("/wss", wss.WsServer())
	r_api.Post("/SetKey", api.SetKey)
	r_api.Post("/HandleKey", api.HandleKey)
	r_api.Post("/GetAccountDetail", api.GetAccountDetail)
	r_api.Post("/Order", api.Order)
	r_api.Post("/EditConfig", api.EditConfig)
	r_api.Post("/SetAccountConfig", api.SetAccountConfig)

	// sys
	s_api := app.Group("/CoinAI/sys")
	s_api.Post("/remove", sys.Remove)
	s_api.Post("/restart", sys.ReStart)
	s_api.Post("/TheOpen", sys.TheOpen)

	// Ping
	app.Use(api.Ping)

	listenHost := mStr.Join(":", config.AppEnv.Port)
	global.Log.Println(mStr.Join(`启动服务: http://127.0.0.1`, listenHost))
	app.Listen(listenHost)
}
