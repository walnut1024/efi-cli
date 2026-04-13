package fund

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/model"
)

const fundBase = "https://fundmobapi.eastmoney.com/FundMNewApi"

// GetHistory returns fund NAV history
func GetHistory(code string, pz int) ([]model.Record, error) {
	params := url.Values{
		"FCODE":         {code},
		"IsShareNet":    {"true"},
		"MobileKey":     {"1"},
		"appType":       {"ttjj"},
		"appVersion":    {"6.2.8"},
		"cToken":        {"1"},
		"deviceid":      {"1"},
		"pageIndex":     {"1"},
		"pageSize":      {fmt.Sprintf("%d", pz)},
		"plat":          {"Iphone"},
		"product":       {"EFund"},
		"serverVersion": {"6.2.8"},
		"uToken":        {"1"},
		"userId":        {"1"},
		"version":       {"6.2.8"},
	}

	data, err := client.DefaultClient.Get(fundBase+"/FundMNHisNetList", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Datas []map[string]interface{} `json:"Datas"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	records := make([]model.Record, 0, len(resp.Datas))
	for _, item := range resp.Datas {
		r := make(model.Record)
		r["code"] = code
		r["date"] = strVal(item, "FSRQ")
		r["nav"] = model.CleanValue(strVal(item, "DWJZ"))
		r["acc_nav"] = model.CleanValue(strVal(item, "LJJZ"))
		r["pct"] = model.CleanValue(strVal(item, "JZZZL"))
		records = append(records, r)
	}
	return records, nil
}

// GetRealtime returns realtime fund estimate
func GetRealtime(codes []string) ([]model.Record, error) {
	params := url.Values{
		"pageIndex":  {"1"},
		"pageSize":   {"300000"},
		"Sort":       {""},
		"Fcodes":     {strings.Join(codes, ",")},
		"SortColumn": {""},
		"IsShowSE":   {"false"},
		"P":          {"F"},
		"deviceid":   {"3EA024C2-7F22-408B-95E4-383D38160FB3"},
		"plat":       {"Iphone"},
		"product":    {"EFund"},
		"version":    {"6.2.8"},
	}

	data, err := client.DefaultClient.Get(fundBase+"/FundMNFInfo", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Datas []map[string]interface{} `json:"Datas"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	records := make([]model.Record, 0, len(resp.Datas))
	for _, item := range resp.Datas {
		r := make(model.Record)
		r["code"] = strVal(item, "FCODE")
		r["name"] = strVal(item, "SHORTNAME")
		r["nav"] = model.CleanValue(strVal(item, "ACCNAV"))
		r["nav_date"] = strVal(item, "PDATE")
		r["estimate_time"] = strVal(item, "GZTIME")
		r["estimate_pct"] = model.CleanValue(strVal(item, "GSZZL"))
		records = append(records, r)
	}
	return records, nil
}

// GetBaseInfo returns fund basic information
func GetBaseInfo(code string) (model.Record, error) {
	params := url.Values{
		"FCODE":    {code},
		"deviceid": {"3EA024C2-7F22-408B-95E4-383D38160FB3"},
		"plat":     {"Iphone"},
		"product":  {"EFund"},
		"version":  {"6.3.8"},
	}

	data, err := client.DefaultClient.Get(fundBase+"/FundMNNBasicInformation", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Datas map[string]interface{} `json:"Datas"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if resp.Datas == nil {
		return nil, fmt.Errorf("fund code not found: %s", code)
	}

	r := make(model.Record)
	r["code"] = strVal(resp.Datas, "FCODE")
	r["name"] = strVal(resp.Datas, "SHORTNAME")
	r["estab_date"] = strVal(resp.Datas, "ESTABDATE")
	r["pct"] = model.CleanValue(strVal(resp.Datas, "RZDF"))
	r["nav"] = model.CleanValue(strVal(resp.Datas, "DWJZ"))
	r["company"] = strVal(resp.Datas, "JJGS")
	r["nav_date"] = strVal(resp.Datas, "FSRQ")
	r["desc"] = strVal(resp.Datas, "COMMENTS")
	return r, nil
}

// GetInvestPosition returns fund stock holdings
func GetInvestPosition(code, date string) ([]model.Record, error) {
	params := url.Values{
		"FCODE":         {code},
		"appType":       {"ttjj"},
		"deviceid":      {"3EA024C2-7F22-408B-95E4-383D38160FB3"},
		"plat":          {"Iphone"},
		"product":       {"EFund"},
		"serverVersion": {"6.2.8"},
		"version":       {"6.2.8"},
	}
	if date != "" {
		params.Set("DATE", date)
	}

	data, err := client.DefaultClient.Get(fundBase+"/FundMNInverstPosition", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Expansion string          `json:"Expansion"`
		Datas     json.RawMessage `json:"Datas"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	// Datas can be {"fundStocks": [...]} or an array
	var items []map[string]interface{}
	if err := json.Unmarshal(resp.Datas, &items); err != nil {
		var obj map[string]interface{}
		if err2 := json.Unmarshal(resp.Datas, &obj); err2 != nil {
			return nil, fmt.Errorf("decode fund position response: %w", err2)
		}
		if raw, ok := obj["fundStocks"]; ok {
			if arr, ok := raw.([]interface{}); ok {
				for _, v := range arr {
					if m, ok := v.(map[string]interface{}); ok {
						items = append(items, m)
					}
				}
			}
		}
	}

	records := make([]model.Record, 0, len(items))
	for _, item := range items {
		r := make(model.Record)
		r["fund_code"] = code
		r["stock_code"] = strVal(item, "GPDM")
		r["stock_name"] = strVal(item, "GPJC")
		r["pct"] = model.CleanValue(strVal(item, "JZBL"))
		r["change"] = model.CleanValue(strVal(item, "PCTNVCHG"))
		r["action"] = strVal(item, "PCTNVCHGTYPE")
		r["sector"] = strVal(item, "INDEXNAME")
		r["date"] = resp.Expansion
		records = append(records, r)
	}
	return records, nil
}

// GetPeriodChange returns fund period returns
func GetPeriodChange(code string) ([]model.Record, error) {
	params := url.Values{
		"AppVersion": {"6.3.8"},
		"FCODE":      {code},
		"MobileKey":  {"3EA024C2-7F22-408B-95E4-383D38160FB3"},
		"OSVersion":  {"14.3"},
		"deviceid":   {"3EA024C2-7F22-408B-95E4-383D38160FB3"},
		"passportid": {"3061335960830820"},
		"plat":       {"Iphone"},
		"product":    {"EFund"},
		"version":    {"6.3.6"},
	}

	data, err := client.DefaultClient.Get(fundBase+"/FundMNPeriodIncrease", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Datas []map[string]interface{} `json:"Datas"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	titles := []string{"1w", "1m", "3m", "6m", "1y", "2y", "3y", "5y", "ytd", "all"}

	records := make([]model.Record, 0, len(resp.Datas))
	for i, item := range resp.Datas {
		r := make(model.Record)
		r["code"] = code
		r["pct"] = model.CleanValue(strVal(item, "syl"))
		r["avg"] = model.CleanValue(strVal(item, "avg"))
		r["rank"] = model.CleanValue(strVal(item, "rank"))
		r["total"] = model.CleanValue(strVal(item, "sc"))
		if i < len(titles) {
			r["period"] = titles[i]
		}
		records = append(records, r)
	}
	return records, nil
}

// GetAssetAllocation returns fund asset type allocation
func GetAssetAllocation(code, date string) ([]model.Record, error) {
	params := url.Values{
		"FCODE":         {code},
		"OSVersion":     {"14.3"},
		"appVersion":    {"6.3.8"},
		"deviceid":      {"3EA024C2-7F21-408B-95E4-383D38160FB3"},
		"plat":          {"Iphone"},
		"product":       {"EFund"},
		"serverVersion": {"6.3.6"},
		"version":       {"6.3.8"},
	}
	if date != "" {
		params.Set("DATE", date)
	}

	data, err := client.DefaultClient.Get(fundBase+"/FundMNAssetAllocationNew", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Datas []map[string]interface{} `json:"Datas"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	records := make([]model.Record, 0, len(resp.Datas))
	for _, item := range resp.Datas {
		r := make(model.Record)
		r["code"] = code
		r["stock_pct"] = model.CleanValue(strVal(item, "GP"))
		r["bond_pct"] = model.CleanValue(strVal(item, "ZQ"))
		r["cash_pct"] = model.CleanValue(strVal(item, "HB"))
		r["total_size"] = model.CleanValue(strVal(item, "JZC"))
		r["other_pct"] = model.CleanValue(strVal(item, "QT"))
		records = append(records, r)
	}
	return records, nil
}

// sortKeyMap maps user-friendly sort keys to upstream API sc parameter
var sortKeyMap = map[string]string{
	"1w": "rzdf",
	"1m": "1yzf",
	"3m": "3yzf",
	"6m": "6yzf",
	"1y": "1nzf",
	"2y": "2nzf",
	"3y": "3nzf",
	"5y": "3nzf", // fallback to 3y fetch, then enrich with period API
}

// GetFundRank returns fund ranking sorted by performance
func GetFundRank(ft, sortKey string, pn int) ([]model.Record, error) {
	if ft == "" {
		ft = "all"
	}
	sc, ok := sortKeyMap[sortKey]
	if !ok {
		sc = "3nzf"
	}
	if pn <= 0 {
		pn = 50
	}

	// For 5y sort, fetch more candidates to enrich
	fetchN := pn
	if sortKey == "5y" {
		fetchN = 200
	}

	records, err := fetchFundRank(ft, sc, fetchN)
	if err != nil {
		return nil, err
	}

	// Enrich with 5y data from period API
	if sortKey == "5y" {
		enrich5yData(records)
		sort.Slice(records, func(i, j int) bool {
			vi, _ := records[i]["y5"].(float64)
			vj, _ := records[j]["y5"].(float64)
			return vi > vj
		})
		if len(records) > pn {
			records = records[:pn]
		}
	}

	return records, nil
}

func enrich5yData(records []model.Record) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, client.MaxConcurrency)
	var mu sync.Mutex

	for i := range records {
		code, _ := records[i]["code"].(string)
		if code == "" {
			continue
		}
		wg.Add(1)
		go func(idx int, c string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			periods, err := GetPeriodChange(c)
			if err != nil {
				return
			}
			for _, p := range periods {
				if period, ok := p["period"].(string); ok && period == "5y" {
					mu.Lock()
					records[idx]["y5"] = p["pct"]
					mu.Unlock()
					return
				}
			}
		}(i, code)
	}
	wg.Wait()
}

func fetchFundRank(ft, sc string, pn int) ([]model.Record, error) {
	params := url.Values{
		"op":         {"ph"},
		"dt":         {"kf"},
		"ft":         {ft},
		"rs":         {""},
		"gs":         {"0"},
		"sc":         {sc},
		"st":         {"desc"},
		"sd":         {""},
		"ed":         {""},
		"qdii":       {""},
		"tabSubtype": {",,,,,"},
		"pi":         {"1"},
		"pn":         {fmt.Sprintf("%d", pn)},
		"dx":         {"1"},
		"v":          {"0.123"},
	}

	rawURL := "https://fund.eastmoney.com/data/rankhandler.aspx?" + params.Encode()
	req, _ := http.NewRequest("GET", rawURL, nil)
	req.Header.Set("Referer", "https://fund.eastmoney.com/")
	req.Header.Set("User-Agent", client.UserAgent)
	resp, err := client.DefaultClient.HTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	// Parse JS format: var rankData = {datas:["...","..."],allRecords:N,...}
	body := string(data)
	start := strings.Index(body, "datas:[")
	if start < 0 {
		return nil, fmt.Errorf("unexpected rank API response")
	}
	body = body[start+7:] // skip 'datas:['
	end := strings.Index(body, "],allRecords")
	if end < 0 {
		return nil, fmt.Errorf("unexpected rank API response")
	}
	body = body[:end]

	// Split entries: "field1,field2,...","field1,field2,..."
	items := strings.Split(body, "\",\"")
	if len(items) == 0 {
		return []model.Record{}, nil
	}
	items[0] = strings.TrimPrefix(items[0], "\"")
	items[len(items)-1] = strings.TrimSuffix(items[len(items)-1], "\"")

	records := make([]model.Record, 0, len(items))
	for _, item := range items {
		f := strings.Split(item, ",")
		if len(f) < 17 {
			continue
		}
		r := make(model.Record)
		r["code"] = f[0]
		r["name"] = f[1]
		r["date"] = f[3]
		r["nav"] = model.CleanValue(f[4])
		r["acc_nav"] = model.CleanValue(f[5])
		r["day"] = model.CleanValue(f[6])
		r["week"] = model.CleanValue(f[7])
		r["month"] = model.CleanValue(f[8])
		r["m3"] = model.CleanValue(f[9])
		r["m6"] = model.CleanValue(f[10])
		r["y1"] = model.CleanValue(f[11])
		r["y2"] = model.CleanValue(f[12])
		r["y3"] = model.CleanValue(f[13])
		r["total"] = model.CleanValue(f[15])
		r["inception_date"] = f[16]
		records = append(records, r)
	}
	return records, nil
}

// GetFundScreen screens funds by type, size, and return thresholds
func GetFundScreen(ft string, minSize, maxSize, y1Min, y3Min, y5Min float64, sortKey string, limit int) ([]model.Record, error) {
	if ft == "" {
		ft = "all"
	}
	sc, ok := sortKeyMap[sortKey]
	if !ok {
		sc = "3nzf"
	}

	// Fetch a large pool to filter from
	records, err := fetchFundRank(ft, sc, 500)
	if err != nil {
		return nil, err
	}

	// Enrich with 5y data if needed for sorting or filtering
	if sortKey == "5y" || y5Min > 0 {
		enrich5yData(records)
		if sortKey == "5y" {
			sort.Slice(records, func(i, j int) bool {
				vi, _ := records[i]["y5"].(float64)
				vj, _ := records[j]["y5"].(float64)
				return vi > vj
			})
		}
	}

	// Filter
	filtered := make([]model.Record, 0, len(records))
	for _, r := range records {
		if minSize > 0 {
			v := floatVal(r, "total")
			if v < minSize*1e8 {
				continue
			}
		}
		if maxSize > 0 {
			v := floatVal(r, "total")
			if v > maxSize*1e8 {
				continue
			}
		}
		if y1Min > 0 {
			v := floatVal(r, "y1")
			if v < y1Min {
				continue
			}
		}
		if y3Min > 0 {
			v := floatVal(r, "y3")
			if v < y3Min {
				continue
			}
		}
		if y5Min > 0 {
			v, _ := r["y5"].(float64)
			if v < y5Min {
				continue
			}
		}
		filtered = append(filtered, r)
	}

	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}
	return filtered, nil
}

func floatVal(r model.Record, key string) float64 {
	v, ok := r[key]
	if !ok {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case int64:
		return float64(val)
	default:
		return 0
	}
}

func strVal(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}
