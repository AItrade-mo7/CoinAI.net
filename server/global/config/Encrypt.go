package config

import (
	"time"

	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mStr"
)

var SecretKey = mEncrypt.MD5("AITrade.net from mo7cc")

func Encrypt(msg string) string {
	now := time.Now().Unix() / 30 // 30秒一验证

	EnStr := ""
	for i := -2; i < 3; i++ {
		timestamp := now + int64(i)
		s := mEncrypt.Sha256(
			mStr.Join(msg, "mo7", timestamp),
			SecretKey)
		EnStr += s
	}

	return EnStr
}

func ClientEncrypt(msg string) string {
	now := time.Now().Unix() / 30
	return mEncrypt.Sha256(
		mStr.Join(msg, "mo7", now),
		SecretKey)
}
