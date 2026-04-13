package stock

import (
	"strings"

	"github.com/walnut1024/efi-cli/internal/client"
	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

var quoteFields = strings.Join(model.QuoteFieldKeys, ",") + ",f13,f124,f297"

type targetResolver func(string) (string, error)

func resolveQuoteSecIDs(codes []string) ([]string, error) {
	secids := make([]string, 0, len(codes))
	for _, code := range codes {
		secid, err := client.ResolveQuoteID(code)
		if err != nil {
			return nil, err
		}
		secids = append(secids, secid)
	}
	return secids, nil
}

func fetchQuoteRecords(codes []string) ([]model.Record, error) {
	secids, err := resolveQuoteSecIDs(codes)
	if err != nil {
		return nil, err
	}
	items, err := client.FetchQuotesBySecIDs(strings.Join(secids, ","), quoteFields)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", strings.Join(codes, ","), "stock_quote", err, nil)
	}
	return model.ParseQuoteRecords(items), nil
}

func fetchQuoteRaw(codes []string) ([]byte, error) {
	secids, err := resolveQuoteSecIDs(codes)
	if err != nil {
		return nil, err
	}
	_, raw, err := client.FetchQuotesBySecIDsRaw(strings.Join(secids, ","), quoteFields)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", strings.Join(codes, ","), "stock_quote", err, nil)
	}
	return raw, nil
}

func fetchResolvedKlines(target, beg, end string, klt, fqt int, op string, resolver targetResolver) ([]model.Record, error) {
	secid, err := resolver(target)
	if err != nil {
		return nil, err
	}
	resp, err := client.FetchKlines(secid, beg, end, klt, fqt)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", target, op, err, nil)
	}
	return model.ParseKlineRecords(resp.Name, client.SecIDCode(secid), resp.Klines), nil
}

func fetchResolvedKlinesRaw(target, beg, end string, klt, fqt int, op string, resolver targetResolver) ([]byte, error) {
	secid, err := resolver(target)
	if err != nil {
		return nil, err
	}
	_, raw, err := client.FetchKlinesRaw(secid, beg, end, klt, fqt)
	if err != nil {
		return nil, cliutil.NewError(cliutil.ErrUpstream, "上游接口异常", target, op, err, nil)
	}
	return raw, nil
}

func resolveMarketArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}
