package account

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mStr"
)

func LogErr(ApiKey dbType.OkxKeyType, sum ...any) {
	str := fmt.Sprintf("okxApi.LogErr: %+v", sum)
	global.Log.Println(str)

	message := ""
	if len(sum) > 0 {
		message = mStr.ToStr(sum[0])
	}
	content := mJson.Format(sum)

	err := taskPush.SysEmail(taskPush.SysEmailOpt{
		From:        config.AppEnv.SysName,
		To:          []string{ApiKey.UserID},
		Subject:     "交易所接口报错",
		Title:       mStr.Join(ApiKey.Name, " 您在 ", config.AppEnv.SysName, " 的ApiKey出现了错误, 请检查交易所账户设置"),
		Message:     message,
		Content:     content,
		Description: "交易所接口报错",
	})
	global.Log.Println("邮件已发送", err)
}
