package taskPush

import (
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mFetch"
	"github.com/EasyGolang/goTools/mStr"
)

type RequestOpt struct {
	Origin string
	Path   string
	Data   []byte
}

func Request(opt RequestOpt) (resData []byte, resErr error) {
	UserAgent := config.SysName
	Path := opt.Path

	Data := opt.Data
	enData := mEncrypt.MD5(mStr.ToStr(Data))

	fetch := mFetch.NewHttp(mFetch.HttpOpt{
		Origin: opt.Origin,
		Path:   Path,
		Data:   Data,
		Header: map[string]string{
			"Auth-Encrypt": config.ClientEncrypt(Path + UserAgent + enData),
			"User-Agent":   UserAgent,
		},
	})

	return fetch.Post()
}
