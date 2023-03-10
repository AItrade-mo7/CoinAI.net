package taskPush

import (
	"fmt"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mRes"
	"github.com/EasyGolang/goTools/mTask"
	"github.com/EasyGolang/goTools/mTime"
	jsoniter "github.com/json-iterator/go"
)

type NewOpt struct {
	TaskType    string         // 任务类型
	Content     map[string]any // 任务内容
	Description string         // 任务描述
}

func New(opt NewOpt) error {
	now := mTime.GetTime()

	NewTaskData := mTask.TaskType{
		TaskID:        mEncrypt.GetUUID(),
		TaskType:      opt.TaskType,
		Content:       opt.Content,
		Source:        config.SysName,
		Description:   opt.Description, // 任务描述
		CreateTime:    now.TimeUnix,
		CreateTimeStr: now.TimeStr,
	}

	// 发送任务
	resData, err := Request(RequestOpt{
		Origin: config.SysEnv.MessageBaseUrl,
		Path:   "/api/async/InsertTaskQueue",
		Data:   mJson.ToJson(NewTaskData),
	})
	if err != nil {
		return err
	}

	var resObj mRes.ResType
	jsoniter.Unmarshal(resData, &resObj)

	if resObj.Code < 0 {
		return fmt.Errorf(resObj.Msg)
	}
	return err
}
