package binanceApi

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	jsoniter "github.com/json-iterator/go"
)

type SymbolType struct {
	Symbol                     string
	Status                     string
	BaseAsset                  string
	BaseAssetPrecision         int
	QuoteAsset                 string
	QuotePrecision             int
	QuoteAssetPrecision        int
	BaseCommissionPrecision    int
	QuoteCommissionPrecision   int
	OrderTypes                 []string
	IcebergAllowed             bool
	OcoAllowed                 bool
	QuoteOrderQtyMarketAllowed bool
	AllowTrailingStop          bool
	CancelReplaceAllowed       bool
	IsSpotTradingAllowed       bool
	IsMarginTradingAllowed     bool
	Filters                    []struct {
		FilterType            string
		MinPrice              string
		MaxPrice              string
		TickSize              string
		MultiplierUp          string
		MultiplierDown        string
		AvgPriceMins          int
		MinQty                string
		MaxQty                string
		StepSize              string
		MinNotional           string
		ApplyToMarket         bool
		Limit                 int
		MinTrailingAboveDelta int
		MaxTrailingAboveDelta int
		MinTrailingBelowDelta int
		MaxTrailingBelowDelta int
		MaxNumOrders          int
		MaxNumAlgoOrders      int
	}
	Permissions []string
}

type InstType struct {
	Timezone   string
	ServerTime int64
	RateLimits []struct {
		RateLimitType string
		Interval      string
		IntervalNum   int
		Limit         int
	}
	ExchangeFilters []interface{}
	Symbols         []SymbolType
}

func GetInst() (InstList []SymbolType) {
	Kdata_file := mStr.Join(config.Dir.JsonData, "/B-Inst", ".json")

	resData, err := mOKX.FetchBinance(mOKX.FetchBinanceOpt{
		Path:          "/api/v3/exchangeInfo",
		Method:        "get",
		LocalJsonPath: Kdata_file,
	})
	if err != nil {
		global.LogErr("binanceApi.GetInst Err", err)
	}
	var result InstType
	jsoniter.Unmarshal(resData, &result)

	if len(result.Symbols) < 10 {
		global.LogErr("binanceApi.GetInst 长度不足", len(result.Symbols))
		return
	}

	InstList = []SymbolType{}
	for _, val := range result.Symbols {
		if val.QuoteAsset == "USDT" && val.Status == "TRADING" {
			InstList = append(InstList, val)
		}
	}

	mFile.Write(Kdata_file, mStr.ToStr(resData))
	return
}
