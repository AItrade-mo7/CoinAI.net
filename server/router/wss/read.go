package wss

import (
	"strings"
	"time"

	"Hunter.net/server/global/config"
	"Hunter.net/server/router/result"
	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mRes"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

func Read(msg []byte) mRes.ResType {
	MsgStr := mStr.ToStr(msg)

	if strings.ToLower(MsgStr) == "ping" {
		return result.Succeed.WithData("pong")
	}

	findToken := strings.Contains(MsgStr, `"Token":"`)
	findEncrypt := strings.Contains(MsgStr, `"Encrypt":`)

	if findToken || findEncrypt {
		return verifyCode(msg)
	}

	return result.Succeed.WithData(msg)
}

func verifyCode(data []byte) mRes.ResType {
	var Auth struct {
		Token   string `json:"Token"`
		Encrypt string `json:"Encrypt"`
	}
	err := jsoniter.Unmarshal(data, &Auth)
	if err != nil {
		return result.Fail.WithData("数据格式化失败")
	}

	shaStr := config.Encrypt("/wss")

	if len(Auth.Encrypt) < 1 {
		return result.ErrAuth.WithData("缺少授权码")
	}

	isFind := strings.Contains(shaStr, Auth.Encrypt)

	if !isFind {
		return result.ErrAuth.WithData("授权验证失败")
	}

	if len(Auth.Token) < 1 {
		return result.ErrToken.WithData("缺少Token")
	}

	Claims, TokenOK := mEncrypt.ParseToken(Auth.Token, config.SecretKey)
	if !TokenOK {
		return result.ErrToken.WithData("Token验证失败")
	}

	UserID := Claims.Message
	if len(UserID) != 32 {
		return result.ErrToken.WithData("Token解析失败")
	}

	ExpiresAt := Claims.StandardClaims.ExpiresAt
	now := time.Now().Unix()

	if ExpiresAt-now < 0 {
		return result.ErrToken.WithData("Token过期")
	}

	return result.RightLogin.WithData("")
}
