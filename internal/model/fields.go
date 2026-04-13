package model

// QuoteFields maps API field IDs to output key names for realtime quotes
var QuoteFields = map[string]string{
	"f12":  "code",
	"f14":  "name",
	"f3":   "pct",
	"f2":   "price",
	"f15":  "high",
	"f16":  "low",
	"f17":  "open",
	"f4":   "chg",
	"f8":   "turnover",
	"f10":  "vol_ratio",
	"f9":   "pe",
	"f5":   "vol",
	"f6":   "amount",
	"f18":  "pre_close",
	"f20":  "mkt_cap",
	"f21":  "float_cap",
	"f13":  "mkt_num",
	"f124": "update_ts",
	"f297": "latest_trade_date",
}

// KlineFields maps K-line API field IDs to output key names
var KlineFields = map[string]string{
	"f51": "date",
	"f52": "open",
	"f53": "close",
	"f54": "high",
	"f55": "low",
	"f56": "vol",
	"f57": "amount",
	"f58": "amplitude",
	"f59": "pct",
	"f60": "chg",
	"f61": "turnover",
}

// KlineFieldKeys is the ordered list of K-line API field IDs
var KlineFieldKeys = []string{"f51", "f52", "f53", "f54", "f55", "f56", "f57", "f58", "f59", "f60", "f61"}

// QuoteFieldKeys is the ordered list of quote API field IDs
var QuoteFieldKeys = []string{"f12", "f14", "f3", "f2", "f15", "f16", "f17", "f4", "f8", "f10", "f9", "f5", "f6", "f18", "f20", "f21", "f13"}

// BillFields maps capital flow API field IDs
var BillFields = map[string]string{
	"f51": "date",
	"f52": "main_net",
	"f53": "small_net",
	"f54": "med_net",
	"f55": "big_net",
	"f56": "huge_net",
	"f57": "main_pct",
	"f58": "small_pct",
	"f59": "med_pct",
	"f60": "big_pct",
	"f61": "huge_pct",
	"f62": "close",
	"f63": "pct",
}

// BillFieldKeys is the ordered list of bill API field IDs
var BillFieldKeys = []string{"f51", "f52", "f53", "f54", "f55", "f56", "f57", "f58", "f59", "f60", "f61", "f62", "f63"}

// BaseInfoFields maps base info API field IDs
var BaseInfoFields = map[string]string{
	"f57":  "code",
	"f58":  "name",
	"f162": "pe",
	"f167": "pb",
	"f127": "industry",
	"f116": "mkt_cap",
	"f117": "float_cap",
	"f198": "board_code",
	"f173": "roe",
	"f187": "net_margin",
	"f105": "net_profit",
	"f186": "gross_margin",
}

// DefaultQuoteOutputFields is the default fields for quote output
var DefaultQuoteOutputFields = []string{"code", "name", "pct", "price", "high", "low", "open", "chg", "turnover", "pe", "vol", "amount", "mkt"}

// DefaultKlineOutputFields is the default fields for K-line output
var DefaultKlineOutputFields = []string{"name", "code", "date", "open", "close", "high", "low", "vol", "amount", "pct", "chg", "turnover"}

// DefaultBillOutputFields is the default fields for bill output
var DefaultBillOutputFields = []string{"date", "main_net", "small_net", "med_net", "big_net", "huge_net", "close", "pct"}
