package result

import "github.com/EasyGolang/goTools/mRes"

var (
	Succeed    = mRes.Response(1, "Succeed") // 通用成功
	RightLogin = mRes.Response(2, "登录成功")    // 登录成功

	Fail               = mRes.Response(-1, "Fail") // 通用错误
	ErrAuth            = mRes.Response(-2, "授权验证失败")
	ErrToken           = mRes.Response(-3, "Token验证失败")
	ErrPassword        = mRes.Response(-4, "密码错误")
	ErrDB              = mRes.Response(-5, "数据库出错")
	ErrPower           = mRes.Response(-6, "当前用户无权限")
	ErrAccount         = mRes.Response(-7, "该账号不存在")
	ErrRmUser          = mRes.Response(-8, "注册失败")
	ErrLogin           = mRes.Response(-9, "登录失败")
	ErrUpload          = mRes.Response(-10, "上传失败")
	ErrEmail           = mRes.Response(-11, "邮件发送失败")
	ErrEmailCode       = mRes.Response(-12, "验证码错误")
	ErrAccountRepeat   = mRes.Response(-13, "该账号已存在")
	ErrAddOkxKey       = mRes.Response(-14, "密钥创建失败")
	ErrAddHunterServer = mRes.Response(-15, "服务创建失败")
)
