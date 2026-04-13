package stock

import (
	"net/url"

	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

var commentFields = map[string]string{
	"SECURITY_CODE": "code", "SECURITY_NAME_ABBR": "name", "CLOSE_PRICE": "price",
	"CHANGE_RATE": "pct", "TURNOVERRATE": "turnover", "PE_DYNAMIC": "pe",
	"PRIME_COST": "main_cost", "ORG_PARTICIPATE": "institution_participation",
	"TOTALSCORE": "score", "RATIO": "rising", "RANK": "rank",
	"FOCUS": "attention", "TRADE_DATE": "date",
}

var commentSchema = cliutil.CommandSchema{
	Command: "stock comment",
	Entity:  "stock_comment",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"}, "sort": true,
	},
	DefaultFields: []string{"code", "name", "price", "pct", "pe", "score", "rank", "attention", "main_cost"},
	Fields: []cliutil.FieldSchema{
		{Name: "code", Type: "string", Desc: "证券代码"},
		{Name: "name", Type: "string", Desc: "证券名称"},
		{Name: "price", Type: "number", Desc: "收盘价"},
		{Name: "pct", Type: "number", Desc: "涨跌幅"},
		{Name: "turnover", Type: "number", Desc: "换手率"},
		{Name: "pe", Type: "number", Desc: "市盈率"},
		{Name: "main_cost", Type: "number", Desc: "主力成本"},
		{Name: "institution_participation", Type: "number", Desc: "机构参与度"},
		{Name: "score", Type: "number", Desc: "综合评分"},
		{Name: "rising", Type: "number", Desc: "上升比例"},
		{Name: "rank", Type: "number", Desc: "排名"},
		{Name: "attention", Type: "number", Desc: "关注度"},
		{Name: "date", Type: "string", Desc: "日期"},
	},
}

var commentInstitutionSchema = cliutil.CommandSchema{
	Command: "stock comment institution",
	Entity:  "stock_comment_institution",
	Supports: map[string]interface{}{
		"format": []string{"csv", "json", "table"},
	},
	DefaultFields: []string{"date", "institution_participation"},
	Fields: []cliutil.FieldSchema{
		{Name: "date", Type: "string", Desc: "日期"},
		{Name: "institution_participation", Type: "number", Desc: "机构参与度"},
	},
}

func GetComment(limit int) ([]model.Record, error) {
	filters := url.Values{
		"sortColumns": {"TRADE_DATE"},
		"sortTypes":   {"-1"},
	}
	data, err := fetchDatacenter("RPT_DMSK_TS_STOCKNEW", filters, limit)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", "", "comment", err, nil)
	}
	return parseDatacenterToRecords(data, commentFields), nil
}

func GetCommentInstitution(code string) ([]model.Record, error) {
	institutionFields := map[string]string{
		"TRADE_DATE": "date", "ORG_PARTICIPATE": "institution_participation",
	}
	filters := url.Values{
		"filter":      {datacenterFilter(code, "SECURITY_CODE")},
		"sortColumns": {"TRADE_DATE"},
		"sortTypes":   {"-1"},
	}
	data, err := fetchDatacenter("RPT_DMSK_TS_STOCKEVALUATE", filters, 30)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", code, "comment_institution", err, nil)
	}
	records := parseDatacenterToRecords(data, institutionFields)
	for i := range records {
		records[i]["code"] = code
	}
	return records, nil
}
