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

//go:embed Reboot.sh
var Reboot string

type RebootParam struct {
	Port string
	Path string
}

//go:embed Shutdown.sh
var Shutdown string

type ShutdownParam struct {
	Port string
	Path string
}
