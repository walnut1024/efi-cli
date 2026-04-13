# Screening and Compare Reference

Use this file when the user asks for side-by-side analysis, candidates, or a recommendation shortlist.

## Compare Workflow

### Stocks

```bash
efi-cli stock quote <left> <right>
efi-cli stock compare <left> <right> --beg YYYYMMDD --metric total_return
efi-cli stock info <left>
efi-cli stock info <right>
```

### Funds

```bash
efi-cli fund period <left>
efi-cli fund period <right>
efi-cli fund asset <left>
efi-cli fund asset <right>
```

Compare on:

- Return
- Quality or valuation
- Volatility / drawdown profile
- Flow or sentiment support
- Fit for the user's horizon

## Screening Workflow

### Stocks

```bash
efi-cli stock screen --pe-min 5 --pe-max 20 --roe-min 15 --limit 30
efi-cli stock screen --market 创业板 --pct-min 3 --sort pct --limit 20
```

### Funds

```bash
efi-cli fund screen --type gp --size-min 10 --1y-min 10 --3y-min 30 --sort 3y --limit 20
efi-cli fund screen --size-min 5 --size-max 50 --format table
```

## How to Present Screened Results

- Narrow to the best few names instead of echoing the whole list
- Explain why each candidate passed
- Also explain what the filter does not capture
- Validate the top candidates with deeper commands before making a recommendation

## Recommendation Framing

When the user asks "which is better" or "worth buying":

- State the answer relative to a goal and horizon
- Show the deciding factors
- Include the biggest reason you could be wrong
- If neither is attractive, say so directly
