package stock

import (
	"net/url"

	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

var marginFields = map[string]string{
	"STATISTICS_DATE": "date", "FIN_BALANCE": "fin_balance", "LOAN_BALANCE": "loan_balance",
	"FIN_BUY_AMT": "fin_buy_amt", "LOAN_SELL_AMT": "loan_sell_amt",
	"INVESTOR_NUM": "investor_num", "TOTAL_GUARANTEE": "total_guarantee",
	"AVG_GUARANTEE_RATIO": "avg_guarantee_ratio", "SCI_CLOSE_PRICE": "index_close",
	"SCI_CHANGE_RATE": "index_pct",
}

var marginSchema = cliutil.CommandSchema{
	Command: "stock margin account",
	Entity:  "stock_margin_account",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"},
	},
	DefaultFields: []string{"date", "fin_balance", "loan_balance", "fin_buy_amt", "loan_sell_amt", "investor_num", "total_guarantee", "avg_guarantee_ratio"},
	Fields: []cliutil.FieldSchema{
		{Name: "date", Type: "string", Desc: "统计日期"},
		{Name: "fin_balance", Type: "number", Desc: "融资余额(亿元)"},
		{Name: "loan_balance", Type: "number", Desc: "融券余额(亿元)"},
		{Name: "fin_buy_amt", Type: "number", Desc: "融资买入额(亿元)"},
		{Name: "loan_sell_amt", Type: "number", Desc: "融券卖出额(亿元)"},
		{Name: "investor_num", Type: "number", Desc: "参与交易投资者数"},
		{Name: "total_guarantee", Type: "number", Desc: "担保物总价值(亿元)"},
		{Name: "avg_guarantee_ratio", Type: "number", Desc: "平均维持担保比例"},
		{Name: "index_close", Type: "number", Desc: "上证收盘"},
		{Name: "index_pct", Type: "number", Desc: "上证涨跌幅"},
	},
}

func GetMarginAccount(limit int) ([]model.Record, error) {
	filters := url.Values{
		"sortColumns": {"STATISTICS_DATE"},
		"sortTypes":   {"-1"},
	}
	data, err := fetchDatacenter("RPTA_WEB_MARGIN_DAILYTRADE", filters, limit)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", "", "margin_account", err, nil)
	}
	return parseDatacenterToRecords(data, marginFields), nil
}
