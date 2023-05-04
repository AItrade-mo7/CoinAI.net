package account

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mStr"
)

func LogErr(ApiKey dbType.OkxKeyType, sum ...any) {
	str := fmt.Sprintf("okxApi.LogErr: %+v", sum)
	global.Log.Println(str)

	message := ""
	if len(sum) > 0 {
		message = mStr.ToStr(sum[0])
	}
	content := mStr.ToStr(sum)

	err := taskPush.SysEmail(taskPush.SysEmailOpt{
		From:        config.AppEnv.SysName,
		To:          []string{ApiKey.UserID},
		Subject:     "交易所接口报错,请及时检查账户以及持仓！",
		Title:       mStr.Join(ApiKey.Name, ", 您在卫星服务: ", config.AppEnv.SysName, " 绑定的 ApiKey 出现了错误, 请及时检查持仓以及设置！"),
		Message:     message,
		Content:     content,
		Description: "交易所接口报错",
	})
	global.Log.Println(ApiKey.Name, "邮件已发送", err)
}
