package reqDataCenter

import (
	"CoinAI.net/server/global/config"
	jsoniter "github.com/json-iterator/go"
)

type PingDataType struct {
	Code int `json:"Code"`
	Data struct {
		IP string `json:"IP"`
	} `json:"Data"`
	Msg string `json:"Msg"`
}

func GetLocalIP() string {
	resData, _ := NewRest(RestOpt{
		Origin: config.Origin,
		Path:   "/ping",
		Method: "GET",
	})

	var PingData PingDataType
	jsoniter.Unmarshal(resData, &PingData)

	return PingData.Data.IP
}
