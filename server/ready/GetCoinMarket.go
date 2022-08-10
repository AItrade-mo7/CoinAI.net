package ready

import (
	"fmt"
	"time"
)

func GetCoinMarket() {
	time.Sleep(time.Second * 2)
	fmt.Println("获取榜单数据")
}
