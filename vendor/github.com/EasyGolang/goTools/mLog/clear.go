package mLog

import (
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/EasyGolang/goTools/mPath"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

type ClearParam struct {
	Path      string
	ClearTime int64 // 毫秒时长，默认一个周
}

func Clear(opt ClearParam) {
	ClearTime := opt.ClearTime
	if ClearTime < mTime.UnixTimeInt64.Minute*60 {
		ClearTime = mTime.UnixTimeInt64.Day * 7
	}

	logPath := "./logs"
	if len(opt.Path) > 1 {
		logPath = opt.Path
	}

	isLogPath := mPath.Exists(logPath)
	if !isLogPath {
		return
	}

	fileInfoList, _ := os.ReadDir(logPath)
	timeNow := mTime.ToUnixMsec(time.Now())

	for i := range fileInfoList {
		name := fileInfoList[i].Name()
		path := logPath + "/" + name

		if mPath.IsFile(path) {
			timeStr := logNameTime(path)
			tm2, err := time.ParseInLocation("06年1月02日15时", timeStr, time.Local)
			if err != nil {
				continue
			}
			fileUnix := mTime.ToUnixMsec(tm2)

			if (timeNow - fileUnix) > ClearTime {
				os.Remove(path)
			}

		}
	}
}

func logNameTime(name string) string {
	starStr := "-T"
	endStr := ".log"
	messagePat := mStr.Join(
		starStr, `(.*?)`, endStr,
	)
	reg := regexp.MustCompile(messagePat)
	strArr := reg.FindAllString(name, -1)
	if len(strArr) > 0 {
		str := strArr[0]
		str = strings.Replace(str, starStr, "", -1)
		str = strings.Replace(str, endStr, "", -1)
		return str
	} else {
		return "19年11月11日11时"
	}
}
