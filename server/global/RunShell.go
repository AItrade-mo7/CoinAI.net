package global

import (
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTikker"
)

// https://blog.csdn.net/raoxiaoya/article/details/109014347

func SysReStart() {
	ShellCont := mStr.Join("source ", config.File.Reboot)

	mTikker.NewTikker(mTikker.TikkerOpt{
		LogPath:      config.Dir.Log,
		ShellContent: ShellCont,
	}).InstPm2().RunToPm2()
}

func SysRemove() {
	ShellCont := mStr.Join("source ", config.File.Shutdown)

	mTikker.NewTikker(mTikker.TikkerOpt{
		LogPath:      config.Dir.Log,
		ShellContent: ShellCont,
	}).InstPm2().RunToPm2()
}
