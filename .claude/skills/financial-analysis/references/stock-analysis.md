# Stock Analysis Reference

Use this file when the request is mainly about an individual stock or a small set of stocks.

## Minimal Stock Workflow

```bash
efi-cli stock quote <code>
efi-cli stock info <code>
efi-cli stock history <code> --limit 60
```

Add only as needed:

```bash
efi-cli stock bill <code> --limit 20
efi-cli stock board <code>
efi-cli stock deal <code> --max 50
efi-cli stock comment institution <code>
efi-cli stock connect hold <code> --start YYYY-MM-DD
```

## What Each Dataset Is Good For

- `stock quote`: current snapshot, valuation, liquidity, basic context
- `stock info`: fundamentals and static company quality metrics
- `stock history`: trend, drawdown, support/resistance, indicator context
- `stock bill`: multi-day institutional or large-order flow clues
- `stock board`: theme and sector attribution
- `stock deal`: intraday tape detail and abnormal trading behavior
- `stock comment institution`: institutional participation trend
- `stock connect hold`: northbound holding change and float ownership trend (note: current data source is per-organization, not stock-level aggregate; historical coverage is limited)

## Interpretation Checklist

### Trend

- Recent closes vs recent highs/lows
- Whether pullbacks are shallow or deep
- If using indicators, treat them as supporting evidence, not the thesis itself

### Volume and Activity

- Rising price with rising volume is healthier than rising price with fading volume
- High `vol_ratio` or high turnover means the stock is attracting attention, but not necessarily for the right reason
- Large `amplitude` means uncertainty or aggressive trading

### Valuation and Quality

- `pe`, `pb`, `roe`, `gross_margin`, `net_margin`, `net_profit`
- Interpret valuation in context of business type and growth expectations
- For cyclical businesses, valuation ratios alone can mislead

### Flow and Participation

- `main_net` matters more when it is sustained over several sessions
- Northbound accumulation can support medium-term narratives better than one-day ranking changes
- Institutional participation is more useful as a trend than as a single reading

## Good Output Shape for Stock Questions

1. Current setup: trend and trading character
2. Valuation/fundamental read
3. Capital-flow or sentiment confirmation
4. Key risks and what would invalidate the view
