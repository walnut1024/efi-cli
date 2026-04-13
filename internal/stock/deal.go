package stock

import (
	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/model"
)

func GetDealDetail(code string, maxCount int) ([]model.Record, error) {
	secid, err := client.ResolveQuoteID(code)
	if err != nil {
		return nil, err
	}
	resp, err := client.FetchDealDetails(secid, maxCount)
	if err != nil {
		return nil, err
	}
	return model.ParseDealRecords(client.SecIDCode(secid), resp.Details, resp.PrePrice), nil
}
