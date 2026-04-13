package futures

import (
	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/model"
)

// GetRealtime returns realtime futures quotes
func GetRealtime() ([]model.Record, error) {
	resp, err := client.FetchQuoteList("m:113,m:114,m:115,m:8,m:142,m:225", "f12,f14,f3,f2,f15,f16,f17,f4,f8,f10,f9,f5,f6,f18,f20,f21,f13", 1, 200)
	if err != nil {
		return nil, err
	}
	return model.ParseQuoteRecords(resp.Diff), nil
}

// GetHistory returns futures K-line history (secid is used directly, e.g. "115.ZCM")
func GetHistory(secid, beg, end string, klt, fqt int) ([]model.Record, error) {
	resp, err := client.FetchKlines(secid, beg, end, klt, fqt)
	if err != nil {
		return nil, err
	}
	return model.ParseKlineRecords(resp.Name, client.SecIDCode(secid), resp.Klines), nil
}

// GetDealDetail returns futures deal detail (secid used directly)
func GetDealDetail(secid string, maxCount int) ([]model.Record, error) {
	resp, err := client.FetchDealDetails(secid, maxCount)
	if err != nil {
		return nil, err
	}
	return model.ParseDealRecords(client.SecIDCode(secid), resp.Details, resp.PrePrice), nil
}
