package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mStr"
)

func StartEmail() {
	Content := mStr.Join(
		"服务已启动: ", config.AppEnv.ServeID,
		`<br /> <a href="https://trade.mo7.cc/CoinServe/CoinAI?id=`,
		config.AppEnv.ServeID,
		`"> https://trade.mo7.cc/CoinServe/CoinAI?id=`,
		config.AppEnv.ServeID,
		`</a> <br />`,
		"用户昵称: ",
		config.MainUser.NickName,
		"<br />",
		"服务名称: ",
		config.AppEnv.SysName,
		"服务版本: ",
		config.AppEnv.SysVersion,
		"<br />",
	)

	err := taskPush.SysEmail(taskPush.SysEmailOpt{
		From:        config.SysName,
		To:          config.NoticeEmail,
		Subject:     "系统启动",
		Title:       config.SysName + " 系统启动",
		Message:     "系统启动",
		Content:     Content,
		Description: "系统启动邮件",
	})
	global.Run.Println("系统启动邮件已发送", err)
}
