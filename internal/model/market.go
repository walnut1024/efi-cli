package model

// MarketType represents different market classifications
type MarketType string

const (
	MarketAStock       MarketType = "AStock"
	MarketIndex        MarketType = "Index"
	MarketHK           MarketType = "HK"
	MarketUSStock      MarketType = "UsStock"
	MarketLSE          MarketType = "LSE"
	MarketLSEIOB       MarketType = "LSEIOB"
	MarketUniversalIdx MarketType = "UniversalIndex"
	MarketSIXSwiss     MarketType = "SIX"
	MarketNEEQ         MarketType = "NEEQ"
	MarketBK           MarketType = "BK"
)

// MarketNumberDict maps market number to name
var MarketNumberDict = map[string]string{
	"0":   "深A",
	"1":   "沪A",
	"8":   "中金所",
	"90":  "板块",
	"105": "美股",
	"106": "美股",
	"107": "美股",
	"113": "上期所",
	"114": "大商所",
	"115": "郑商所",
	"116": "港股",
	"128": "港股",
	"142": "上海能源",
	"155": "英股",
	"225": "广期所",
}

// FSDict maps Chinese market names to eastmoney fs parameter values
var FSDict = map[string]string{
	"bond":    "b:MK0354",
	"可转债":     "b:MK0354",
	"stock":   "m:0 t:6,m:0 t:80,m:1 t:2,m:1 t:23,m:0 t:81 s:2048",
	"沪深A股":    "m:0 t:6,m:0 t:80,m:1 t:2,m:1 t:23",
	"沪深京A股":   "m:0 t:6,m:0 t:80,m:1 t:2,m:1 t:23,m:0 t:81 s:2048",
	"北证A股":    "m:0 t:81 s:2048",
	"北A":      "m:0 t:81 s:2048",
	"futures": "m:113,m:114,m:115,m:8,m:142,m:225",
	"期货":      "m:113,m:114,m:115,m:8,m:142,m:225",
	"上证A股":    "m:1 t:2,m:1 t:23",
	"沪A":      "m:1 t:2,m:1 t:23",
	"深证A股":    "m:0 t:6,m:0 t:80",
	"深A":      "m:0 t:6,m:0 t:80",
	"新股":      "m:0 f:8,m:1 f:8",
	"创业板":     "m:0 t:80",
	"科创板":     "m:1 t:23",
	"沪股通":     "b:BK0707",
	"深股通":     "b:BK0804",
	"风险警示板":   "m:0 f:4,m:1 f:4",
	"两网及退市":   "m:0 s:3",
	"地域板块":    "m:90 t:1 f:!50",
	"行业板块":    "m:90 t:2 f:!50",
	"概念板块":    "m:90 t:3 f:!50",
	"上证系列指数":  "m:1 s:2",
	"深证系列指数":  "m:0 t:5",
	"沪深系列指数":  "m:1 s:2,m:0 t:5",
	"ETF":     "b:MK0021,b:MK0022,b:MK0023,b:MK0024",
	"LOF":     "b:MK0404,b:MK0405,b:MK0406,b:MK0407",
	"美股":      "m:105,m:106,m:107",
	"港股":      "m:128 t:3,m:128 t:4,m:128 t:1,m:128 t:2",
	"英股":      "m:155 t:1,m:155 t:2,m:155 t:3,m:156 t:1,m:156 t:2,m:156 t:5,m:156 t:6,m:156 t:7,m:156 t:8",
	"中概股":     "b:MK0201",
	"中国概念股":   "b:MK0201",
}

// GetMarketName returns the market name for a market number
func GetMarketName(mktNum string) string {
	if name, ok := MarketNumberDict[mktNum]; ok {
		return name
	}
	return ""
}
