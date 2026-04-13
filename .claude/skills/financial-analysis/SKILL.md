---
name: financial-analysis
description: "Analyze Chinese market securities with efi-cli. Use when the user wants stock, fund, sector, market sentiment, capital flow, screening, or portfolio-style analysis based on public market data, especially for A-share/HK/US symbols tracked by efi-cli. Trigger on requests such as 个股分析, 基金分析, 走势, 估值, 持仓, 板块, 资金流向, 涨停, 北向资金, 融资融券, 解禁, 选股, 选基, 对比, 风险复盘."
---

# Financial Analysis with efi-cli

Use this skill to turn `efi-cli` market data into concise, decision-useful analysis. Keep the main skill workflow-focused. Read reference files only when needed.

Start with `references/efi-commands.md` only if you need exact flags or field names.

## Before You Analyze

First classify the request:

- **Object**: stock, fund, sector/theme, market sentiment, capital flow, screening, or comparison
- **Intent**: snapshot, trend analysis, valuation check, peer comparison, holdings review, shortlist generation, or risk review
- **Time window**: today, recent 5-20 trading days, recent 1-3 months, or long-term

Then choose the smallest useful dataset:

- For a normal stock thesis, start with `quote`, `info`, and `history`
- Add `bill` only when money flow matters
- Add `board`, `sector`, `comment`, or `connect` only when the user asks about theme, sentiment, or northbound behavior
- For funds, start with `info`, `period`, `asset`, and optionally `history`
- For screening requests, do not fetch follow-up details for every result; inspect only the top candidates

Avoid defaulting to wide pulls or long date ranges when a tighter query answers the question.

## Core Workflow

### 1. Stock Analysis

Default path:

```bash
efi-cli stock quote <code>
efi-cli stock info <code>
efi-cli stock history <code> --limit 60
```

Add these only when needed:

- Money flow: `efi-cli stock bill <code> --limit 20`
- Sector context: `efi-cli stock board <code>`
- Intraday/tick context: `efi-cli stock deal <code> --max 50`
- Sentiment or institutional participation: `efi-cli stock comment institution <code>`
- Northbound holding trend: `efi-cli stock connect hold <code> --start YYYY-MM-DD`

Analyze in this order:

1. Price trend: recent direction, key highs/lows, whether the move is accelerating or weakening
2. Trading quality: turnover, volume ratio, amplitude, whether price and volume confirm each other
3. Valuation and quality: PE, PB, ROE, margins, profit scale
4. Capital behavior: main net inflow/outflow, northbound changes, institutional participation if available
5. Risks and uncertainty: stretched valuation, weakening flow, high volatility, sector drag

### 2. Fund Analysis

Default path:

```bash
efi-cli fund info <code>
efi-cli fund period <code>
efi-cli fund asset <code>
efi-cli fund history <code> --limit 60
```

When the user asks about holdings quality:

```bash
efi-cli fund position <code>
efi-cli stock quote <top_holding_codes...>
efi-cli stock info <top_holding_code>
```

Analyze in this order:

1. Return profile: recent period returns, rank/total, relative performance vs category average
2. Risk profile: daily NAV stability, equity exposure, fund size, style concentration
3. Holdings quality: concentration, top holdings valuation/quality, whether holdings match stated style
4. Portfolio fit: aggressive, balanced, or defensive characteristics

### 3. Comparison

For stocks, use both quick snapshots and longer comparison:

```bash
efi-cli stock quote <left> <right>
efi-cli stock compare <left> <right> --beg YYYYMMDD --metric total_return
```

For funds, compare `fund period`, `fund asset`, and if needed `fund position`.

Make the comparison explicit across:

- Return
- Valuation or holdings quality
- Risk/volatility
- Capital support or sentiment
- Which type of investor each better fits

### 4. Market, Sector, and Sentiment Analysis

Use these when the user asks about hot themes, market appetite, or short-term trading mood:

- Sector strength: `stock sector list`, `stock sector members`, `stock sector history`
- Limit-up/down and pool behavior: `stock pool zt|dt|zb|strong|sub-new`
- Popularity: `stock hot-rank list`, `stock hot-rank keyword`
- Market leverage and capital mood: `stock margin account`, `stock connect history`, `stock connect realtime`
- Restricted-share pressure: `stock restricted summary`, `detail`, `queue`

Focus on regime clues:

- Is money concentrating in a theme or dispersing?
- Are strong stocks keeping their gains or failing intraday?
- Is northbound or leverage behavior confirming the tape?
- Are upcoming unlocks a likely supply overhang?

### 5. Screening and Shortlists

Use `stock screen` or `fund screen` to generate candidates, then inspect the top few manually.

Good patterns:

```bash
efi-cli stock screen --pe-min 5 --pe-max 20 --roe-min 15 --limit 30
efi-cli stock screen --market 创业板 --pct-min 3 --sort pct --limit 20
efi-cli fund screen --type gp --size-min 10 --1y-min 10 --3y-min 30 --sort 3y --limit 20
```

After screening:

- Do not present the raw shortlist alone
- Validate top candidates with quote/info/history
- Call out why they passed the filter and what the filter missed

## Output Standards

Every answer should prioritize interpretation over field repetition.

- Lead with the conclusion first, then the evidence
- Keep the structure matched to the request; do not dump every data source you queried
- Include both positives and risks when the user asks whether something is worth buying or holding
- Use relative language when a benchmark is missing: say "looks elevated" or "appears reasonable" instead of asserting certainty
- State uncertainty when the data only supports a short-term or partial view

Use this response shape when the user wants a judgment:

1. Bottom line
2. Core evidence
3. Risks / caveats
4. Optional next watchpoints

## Common Mistakes to Avoid

- Do not call something undervalued or overvalued without stating the basis
- Do not infer long-term quality from one day of flow, popularity, or pool data
- Do not confuse fund performance with holdings quality; evaluate both separately
- Do not use every available command by default
- Do not present a screening result as a recommendation without manual validation
- Do not ignore the user's time horizon; a short-term trader and a long-term allocator need different conclusions

## References

Load these on demand:

- `references/efi-commands.md`: full command and field reference
- `references/stock-analysis.md`: stock-specific data paths and interpretation checklist
- `references/fund-analysis.md`: fund workflows, holdings review, and fund-specific caveats
- `references/market-sentiment.md`: sector, pool, northbound, margin, and unlock analysis
- `references/screening-and-compare.md`: shortlist, compare, and recommendation framing
