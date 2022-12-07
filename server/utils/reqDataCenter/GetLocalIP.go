package reqDataCenter

import (
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
		Origin: "https://trade-api.mo7.cc",
		Path:   "/ping",
		Method: "GET",
	})

	var PingData PingDataType
	jsoniter.Unmarshal(resData, &PingData)

	return PingData.Data.IP
}
