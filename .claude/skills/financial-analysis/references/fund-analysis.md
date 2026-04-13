# Fund Analysis Reference

Use this file when the user asks about a fund, fund ranking, fund screening, or whether a fund fits a portfolio.

## Minimal Fund Workflow

```bash
efi-cli fund info <code>
efi-cli fund period <code>
efi-cli fund asset <code>
efi-cli fund history <code> --limit 60
```

When holdings matter:

```bash
efi-cli fund position <code>
efi-cli stock quote <top_holding_codes...>
efi-cli stock info <top_holding_code>
```

## What to Evaluate

### Performance

- `pct` vs `avg`
- `rank` and `total`
- Whether strength is recent only or persistent across 1y, 3y, and 5y

### Risk Profile

- NAV path smoothness from history
- `stock_pct`, `bond_pct`, `cash_pct`
- Fund size from `total_size`

### Holdings Quality

- Concentration in top holdings
- Whether top holdings are high-quality businesses or crowded expensive names
- Whether the portfolio style matches the fund's label and recent returns

### Practical Portfolio Fit

- High `stock_pct` and concentrated holdings imply higher beta
- Lower equity exposure and higher cash/bond weight can be more defensive
- Strong returns with weak holdings quality deserve skepticism

## Fund-Specific Caveats

- Holdings data is periodic, not real-time
- Great recent returns do not mean the fund is currently low-risk
- Do not rate a fund only by short-term ranking
- Separate manager/process quality from the temporary strength of a single theme

## Good Output Shape for Fund Questions

1. Overall judgment
2. Return quality vs peers
3. Risk and allocation profile
4. Holdings quality and concentration
5. Who the fund may suit, and who it may not
