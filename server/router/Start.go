package router

import (
	"os"
	"time"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/router/api"
	"CoinAI.net/server/router/api/order"
	"CoinAI.net/server/router/api/sys"
	"CoinAI.net/server/router/middle"
	"CoinAI.net/server/router/wss"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
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

	// 跨域
	app.Use(cors.New())
	// 限流
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Second,
	}))
	// 日志中间件
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] [${ip}:${port}] ${status} - ${method} ${latency} ${path} \n",
		TimeFormat: "2006-01-02 - 15:04:05",
		Output:     logFile,
	}), middle.Public, compress.New(), favicon.New())

	// CoinAI
	r_api := app.Group("/CoinAI")
	r_api.Get("/config", api.GetConfig)
	r_api.Get("/wss", wss.WsServer())

	// sys
	s_api := app.Group("/CoinAI/sys")
	s_api.Post("/remove", sys.Remove)
	s_api.Post("/restart", sys.ReStart)
	// 认主
	s_api.Post("/dominate", sys.Dominate)

	// order
	o_api := app.Group("/CoinAI/Order")
	o_api.Post("/Buy", order.Buy)
	o_api.Post("/Sell", order.Sell)
	o_api.Post("/Close", order.Close)

	// Ping
	app.Use(api.Ping)

	listenHost := mStr.Join(":", config.AppEnv.Port)
	global.Log.Println(mStr.Join(`启动服务: http://127.0.0.1`, listenHost))
	app.Listen(listenHost)
}
