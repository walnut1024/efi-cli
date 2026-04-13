package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/walnut1024/efi-cli/internal/cliutil"
)

type Quote struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Pinyin      string `json:"pinyin"`
	ID          string `json:"id"`
	JYS         string `json:"jys"`
	MktNum      string `json:"mkt_num"`
	QuoteID     string `json:"quote_id"`
	UnifiedCode string `json:"unified_code"`
	InnerCode   string `json:"inner_code"`
}

type SearchResponse struct {
	Candidates []*Quote
	Raw        []byte
}

var (
	searchCache  = make(map[string][]*Quote)
	cacheFile    string
	cacheModTime time.Time
	cacheMu      sync.RWMutex
)

func init() {
	home, _ := os.UserHomeDir()
	cacheFile = filepath.Join(home, ".efi", "search_cache.json")
	loadCache()
}

func loadCache() {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return
	}
	var raw map[string]json.RawMessage
	if json.Unmarshal(data, &raw) != nil {
		return
	}
	for k, v := range raw {
		var list []*Quote
		if json.Unmarshal(v, &list) == nil {
			searchCache[k] = list
			continue
		}
		var q Quote
		if json.Unmarshal(v, &q) == nil {
			searchCache[k] = []*Quote{&q}
		}
	}
	cacheModTime = time.Now()
}

func saveCache() {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if time.Since(cacheModTime) < 5*time.Minute {
		return
	}
	if err := os.MkdirAll(filepath.Dir(cacheFile), 0755); err != nil {
		return
	}
	data, _ := json.Marshal(searchCache)
	tmpFile := cacheFile + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return
	}
	if err := os.Rename(tmpFile, cacheFile); err != nil {
		return
	}
	cacheModTime = time.Now()
}

func SearchQuotes(keyword string, count int) (*SearchResponse, error) {
	return searchQuotes(keyword, count, true)
}

func SearchQuotesFresh(keyword string, count int) (*SearchResponse, error) {
	return searchQuotes(keyword, count, false)
}

func searchQuotes(keyword string, count int, useCache bool) (*SearchResponse, error) {
	if strings.TrimSpace(keyword) == "" {
		return nil, cliutil.NewError(cliutil.ErrInvalidArg, "参数错误", keyword, "search", nil, nil)
	}
	if count <= 0 {
		count = 10
	}

	if useCache {
		cacheMu.RLock()
		if cached, ok := searchCache[keyword]; ok && len(cached) > 0 {
			cacheMu.RUnlock()
			return &SearchResponse{Candidates: cloneQuotes(cached, count)}, nil
		}
		cacheMu.RUnlock()
	}

	u := "https://searchapi.eastmoney.com/api/suggest/get"
	params := url.Values{
		"input": {keyword},
		"type":  {"14"},
		"token": {"D43BF722C8E33BDC906FB84D85E326E8"},
		"count": {fmt.Sprintf("%d", count)},
	}

	data, err := DefaultClient.Get(u, params)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", keyword, "search", err, nil)
	}

	var resp struct {
		QuotationCodeTable struct {
			Data []map[string]interface{} `json:"Data"`
		} `json:"QuotationCodeTable"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, cliutil.NewError(cliutil.ErrDecode, "上游响应解析异常", keyword, "search", err, nil)
	}

	items := resp.QuotationCodeTable.Data
	quotes := make([]*Quote, 0, len(items))
	for _, item := range items {
		quotes = append(quotes, parseQuote(item))
	}

	cacheMu.Lock()
	searchCache[keyword] = cloneQuotes(quotes, len(quotes))
	cacheMu.Unlock()
	go saveCache()

	return &SearchResponse{Candidates: cloneQuotes(quotes, count), Raw: data}, nil
}

func SearchQuote(keyword string) (*Quote, error) {
	resp, err := SearchQuotes(keyword, 1)
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.Candidates) == 0 {
		return nil, nil
	}
	return resp.Candidates[0], nil
}

func GetQuoteID(code string) (string, error) {
	if code == "" {
		return "", fmt.Errorf("empty code")
	}
	if secid := RuleBasedSecID(code); secid != "" {
		return secid, nil
	}
	quote, err := SearchQuote(code)
	if err != nil {
		return "", err
	}
	if quote == nil {
		return "", nil
	}
	return quote.QuoteID, nil
}

func ResolveQuoteID(code string) (string, error) {
	if strings.TrimSpace(code) == "" {
		return "", cliutil.NewError(cliutil.ErrInvalidArg, "参数错误", code, "resolve_quote_id", nil, nil)
	}
	if secid := RuleBasedSecID(code); secid != "" {
		return secid, nil
	}

	resp, err := SearchQuotes(code, 10)
	if err != nil {
		return "", err
	}
	if resp == nil || len(resp.Candidates) == 0 {
		return "", cliutil.NewError(cliutil.ErrCodeNotFound, "代码未找到", code, "resolve_quote_id", nil, nil)
	}

	exact := findExactCandidate(code, resp.Candidates)
	if exact != nil {
		return exact.QuoteID, nil
	}

	meta := map[string]any{
		"candidates": resp.Candidates,
	}
	return "", cliutil.NewError(cliutil.ErrAmbiguousCode, "匹配到多个标的", code, "resolve_quote_id", nil, meta)
}

func RuleBasedSecID(code string) string {
	if len(code) == 6 && isAllDigit(code) {
		switch {
		case code[0] == '6':
			return "1." + code
		case code[0] == '0', code[0] == '3':
			return "0." + code
		case code[0] == '2':
			return "0." + code
		case code[0] == '5', code[0] == '7', code[0] == '8':
			return "1." + code
		case code[0] == '4', code[0] == '9':
			return "0." + code
		}
	}
	if len(code) == 4 && isAllDigit(code) {
		return "1." + code
	}
	return ""
}

func isAllDigit(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func parseQuote(rawItem map[string]interface{}) *Quote {
	quote := &Quote{}
	if v, ok := rawItem["Code"].(string); ok {
		quote.Code = v
	}
	if v, ok := rawItem["Name"].(string); ok {
		quote.Name = v
	}
	if v, ok := rawItem["Pinyin"].(string); ok {
		quote.Pinyin = v
	}
	if v, ok := rawItem["ID"].(string); ok {
		quote.ID = v
	}
	if v, ok := rawItem["JYS"].(string); ok {
		quote.JYS = v
	}
	if v, ok := rawItem["MktNum"].(string); ok {
		quote.MktNum = v
	}
	if v, ok := rawItem["QuoteID"].(string); ok {
		quote.QuoteID = v
	}
	if v, ok := rawItem["UnifiedCode"].(string); ok {
		quote.UnifiedCode = v
	}
	if v, ok := rawItem["InnerCode"].(string); ok {
		quote.InnerCode = v
	}
	return quote
}

func cloneQuotes(src []*Quote, limit int) []*Quote {
	if limit <= 0 || limit > len(src) {
		limit = len(src)
	}
	dst := make([]*Quote, 0, limit)
	for i := 0; i < limit; i++ {
		if src[i] == nil {
			continue
		}
		q := *src[i]
		dst = append(dst, &q)
	}
	return dst
}

func findExactCandidate(input string, candidates []*Quote) *Quote {
	for _, candidate := range candidates {
		if candidate == nil {
			continue
		}
		if candidate.Code == input || candidate.Name == input || candidate.UnifiedCode == input || candidate.QuoteID == input {
			return candidate
		}
	}
	if len(candidates) == 1 {
		return candidates[0]
	}
	return nil
}
