package taskPush

import (
	"fmt"

	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mTask"
	"github.com/EasyGolang/goTools/mTime"
)

// ==== 系统邮件 ====
type SysEmailOpt struct {
	From           string   // 缺省 AItrade
	To             []string // 缺省 config.SysEmail
	Subject        string
	Title          string
	Message        string
	Content        string
	Description    string
	EntrapmentCode string // 默认防钓鱼码 trade.mo7.cc
}

func SysEmail(opt SysEmailOpt) error {
	if len(opt.EntrapmentCode) < 1 {
		opt.EntrapmentCode = "trade.mo7.cc"
	}

	if len(opt.From) < 1 {
		opt.From = "AItrade"
	}

	if len(opt.To) < 1 {
		opt.To = []string{config.SysEmail}
	}

	Cont := mJson.StructToMap(mTask.SysEmail{
		To:      opt.To,
		From:    opt.From,
		Subject: opt.Subject,
		SendData: mTask.SysEmailParam{
			Title:          opt.Title,
			Message:        opt.Message,
			Content:        opt.Content,
			SysTime:        mTime.GetTime().TimeStr,
			Source:         config.SysName,
			EntrapmentCode: opt.EntrapmentCode,
		},
	})
	err := New(NewOpt{
		TaskType:    "SysEmail",
		Content:     Cont,
		Description: opt.Description,
	})

	return err
}

// === 发送验证码 ====
type CodeEmailOpt struct {
	To             string
	VerifyCode     string
	Action         string
	EntrapmentCode string
}

func CodeEmail(opt CodeEmailOpt) error {
	if len(opt.To) < 1 || len(opt.VerifyCode) < 1 || len(opt.Action) < 1 || len(opt.EntrapmentCode) < 1 {
		return fmt.Errorf("缺少属性 %v", mJson.Format(opt))
	}

	Cont := mJson.StructToMap(mTask.CodeEmail{
		From:    "AItrade",
		To:      opt.To,
		Subject: "请查收您的验证码",
		SendData: mTask.CodeEmailParam{
			VerifyCode:     opt.VerifyCode,
			Action:         opt.Action,
			SysTime:        mTime.GetTime().TimeStr,
			Source:         config.SysName,
			EntrapmentCode: opt.EntrapmentCode,
		},
	})
	err := New(NewOpt{
		TaskType:    "CodeEmail",
		Content:     Cont,
		Description: "验证码邮件",
	})

	return err
}

// 注册成功通知

type RegisterEmailOpt struct {
	To             string
	Password       string
	EntrapmentCode string // 缺省值 "trade.mo7.cc"
}

func RegisterEmail(opt RegisterEmailOpt) error {
	if len(opt.To) < 1 || len(opt.Password) < 1 || len(opt.EntrapmentCode) < 1 {
		return fmt.Errorf("缺少属性 %v", mJson.Format(opt))
	}

	Cont := mJson.StructToMap(mTask.RegisterSucceedEmail{
		From:    "AItrade",
		To:      opt.To,
		Subject: "注册成功！",
		SendData: mTask.RegisterSucceedEmailParam{
			Password:       opt.Password,
			SysTime:        mTime.UnixFormat(mTime.GetUnixInt64()),
			Source:         config.SysName,
			EntrapmentCode: opt.EntrapmentCode,
		},
	})

	err := New(NewOpt{
		TaskType:    "RegisterSucceedEmail",
		Content:     Cont,
		Description: "注册成功邮件通知",
	})

	return err
}
