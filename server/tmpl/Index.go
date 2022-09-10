package tmpl

import (
	_ "embed"
	"time"
)

//go:embed email-sys.html
var SysEmail string

type SysParam struct {
	Message      string
	SysTime      time.Time
	SecurityCode string
	NickName     string
}

//go:embed Start.html
var StartSlice string

type StartSliceParam struct {
	CoinServeID string
}
