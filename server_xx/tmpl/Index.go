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
}
