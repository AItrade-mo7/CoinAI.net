package ready

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mStr"
)

func StartEmail() {
	layout := `
<div>
	<span class="label">服务名称:</span>
	<span class="val">${SysName}</span>
</div>
<div>
	<span class="label">服务地址:</span>
	<a class="val" href="http://${ServeID}" target="_blank">
		${ServeID}
	</a>
</div>
<div>
	<span class="label">管理界面:</span>
	<a class="val" href="https://trade.mo7.cc/SatelliteServe/CoinAI?id=${ServeID}" target="_blank">
		https://trade.mo7.cc/SatelliteServe/CoinAI?id=${ServeID}
	</a>
</div>
<div>
	<span class="label">服务版本:</span>
	<span class="val">${SysVersion}</span>
</div>
<div>
	<span class="label">所属用户:</span>
	<span class="val">${NickName}</span>
</div>
`

	Content := mStr.Temp(layout, map[string]string{
		"SysName":    config.AppEnv.SysName,
		"ServeID":    config.AppEnv.ServeID,
		"SysVersion": config.AppEnv.SysVersion,
		"NickName":   config.MainUser.NickName,
	})

	err := taskPush.SysEmail(taskPush.SysEmailOpt{
		From:        config.SysName,
		To:          config.NoticeEmail,
		Subject:     "系统启动",
		Title:       config.SysName + " 服务启动",
		Content:     Content,
		Description: "系统启动邮件",
	})
	global.Run.Println("系统启动邮件已发送", err)
}
