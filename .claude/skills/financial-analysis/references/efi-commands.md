# efi-cli CLI Command Reference

## Stock Commands

### `efi-cli stock quote <code1> [code2] ...`
Realtime quote for one or more stocks. Codes are 6-digit (e.g., 600519). Supports `--sort`.

**Fields**: code, name, price, pct, chg, open, high, low, pre_close, vol, amount, turnover, pe, pb, roe, mkt_cap, float_cap, vol_ratio, industry, board_name, secid, mkt_num, update_ts, latest_trade_date

### `efi-cli stock info <code>`
Base info (fundamentals) for a single stock.

**Fields**: code, name, industry, board_name, board_code, price, pct, pe, pb, roe, net_profit, gross_margin, net_margin, mkt_cap, float_cap

### `efi-cli stock history <code> [--klt N] [--fqt N] [--beg YYYYMMDD] [--end YYYYMMDD] [--indicators ...] [--stats ...] [--summary]`
K-line history. Default: daily, forward-adjusted.

**Fields**: name, code, date, open, close, high, low, vol, amount, pct, chg, turnover, amp

**klt values**: 1/5/15/30/60 (minutes), 101=day, 102=week, 103=month
**fqt values**: 0=none, 1=forward, 2=backward
**indicators**: `ma:5,10,20`, `ema:12,26`, `macd`, `rsi:14`, `boll:20,2`
**stats**: `total_return`, `cumulative_pct`, `amplitude_avg`, `high`, `low`, `max_drawdown`, `start_close`, `end_close`, `start_date`, `end_date`, `bars`

### `efi-cli stock realtime [market]`
Market-wide realtime quotes. Default: 沪深京A股. Supports `--sort`.

**Markets**: 沪深A股, 沪深京A股, 创业板, 科创板, 北证A股, etc.

### `efi-cli stock compare <left> <right> [--beg YYYYMMDD] [--end YYYYMMDD] [--metric ...] [--align-date]`
Compare two stocks or indices. Supports Chinese index names.

**metric options**: total_return, max_drawdown, high, low, amplitude_avg, start_close, end_close, bars

### `efi-cli stock bill <code> [--limit N]`
Historical daily capital flow.

**Fields**: date, close, pct, main_net, main_pct, huge_net, huge_pct, big_net, big_pct, med_net, med_pct, small_net, small_pct

### `efi-cli stock todaybill <code>`
Intraday minute-level capital flow.

### `efi-cli stock deal <code> [--max N]`
Latest trading day tick-by-tick deals. Default --max: 10000.

**Fields**: code, time, price, vol, num, pre_close

### `efi-cli stock board <code>`
Sector/board membership for a stock.

**Fields**: board_code, board_name, board_pct, code, name

### `efi-cli stock members <board_code_or_index>`
Get constituent stocks of a board/index. Supports Chinese index names.

### `efi-cli stock billboard [--start YYYY-MM-DD] [--end YYYY-MM-DD]`
Dragon-Tiger list (龙虎榜) data.

### `efi-cli stock holders [--date YYYY-MM-DD]`
Shareholder number changes.

**Fields**: code, name, end_date, holder_num, holder_change, holder_change_pct, avg_hold_num, avg_market_cap, total_shares, total_cap

### `efi-cli stock performance [--date YYYY-MM-DD]`
Company quarterly performance reports.

### `efi-cli stock screen [--market ...] [--pe-min] [--pe-max] [--pb-min] [--pb-max] [--roe-min] [--cap-min] [--cap-max] [--pct-min] [--pct-max] [--limit N]`
Screen stocks by PE, PB, ROE, market cap, daily change.

**Parameters**:
- `--market`: 沪深京A股 (default), 沪深A股, 创业板, 科创板, etc.
- `--pe-min/--pe-max`: PE ratio range
- `--pb-min/--pb-max`: PB ratio range
- `--roe-min`: Min ROE (%)
- `--cap-min/--cap-max`: Market cap range (亿元)
- `--pct-min/--pct-max`: Daily change range (%)

### `efi-cli stock finance <code> [--type income|balance|cashflow]`
Financial statements. Default type: income.

### `efi-cli stock dividend <code>`
Dividend history.

---

## Pool Commands (涨跌停池)

### `efi-cli stock pool zt [--date YYYYMMDD]`
Today's limit-up stocks. Supports `--sort`.

**Fields**: code, name, pct, price, amount, float_cap, total_cap, turnover, seal_amount, first_seal_time, last_seal_time, break_count, consecutive, industry

### `efi-cli stock pool zt-prev [--date YYYYMMDD]`
Previous day limit-up stocks (today's performance). Supports `--sort`.

**Fields**: code, name, pct, price, amount, float_cap, total_cap, turnover, industry

### `efi-cli stock pool dt [--date YYYYMMDD]`
Limit-down stocks. Supports `--sort`.

**Fields**: code, name, pct, price, amount, float_cap, total_cap, pe, turnover, seal_amount, last_seal_time, break_count, consecutive, industry

### `efi-cli stock pool zb [--date YYYYMMDD]`
Broken limit stocks (touched limit but failed to hold). Supports `--sort`.

**Fields**: code, name, pct, price, amount, float_cap, total_cap, turnover, industry

### `efi-cli stock pool strong [--date YYYYMMDD]`
Strong stocks pool. Supports `--sort`.

**Fields**: code, name, pct, price, amount, float_cap, total_cap, turnover, industry

### `efi-cli stock pool sub-new [--date YYYYMMDD]`
Sub-new (recently listed) stocks pool. Supports `--sort`.

**Fields**: code, name, pct, price, amount, float_cap, total_cap, turnover, industry

---

## Sector Commands (板块深数据)

### `efi-cli stock sector list [--type industry|concept]`
Sector list with realtime quotes. Default type: concept. Supports `--sort`.

**Fields**: code, name, price, pct, chg, total_cap, turnover, up_count, down_count, lead_stock, lead_code

### `efi-cli stock sector members <code> [--limit N]`
Sector constituent stocks. Supports `--sort`.

**Fields**: code, name, price, pct, chg, vol, amt, amplitude, high, low, open, pre_close, turnover, pe, pb

### `efi-cli stock sector history <code> [--klt N] [--fqt N] [--beg YYYYMMDD] [--end YYYYMMDD]`
Sector K-line history.

### `efi-cli stock sector quote <code>`
Sector realtime quote.

---

## Stock Connect Commands (沪深港通)

### `efi-cli stock connect summary`
Stock Connect quota summary.

**Fields**: date, mutual_type_name, direction, board_type, index_name

### `efi-cli stock connect history [--type north|sh|sz|south|hshk|szhk] [--limit N]`
Historical fund flow. Default type: north.

**Fields**: date, fund_inflow, net_deal_amt, buy_amt, sell_amt, hold_market_cap, quota_balance, lead_stock, lead_code

### `efi-cli stock connect realtime`
Minute-level intraday fund flow.

**Fields**: time, sh_connect, sz_connect, north_total

### `efi-cli stock connect hold-rank [--type north|sh|sz] [--indicator 1|3|5|10|M|Q|Y] [--date YYYY-MM-DD] [--limit N]`
Holdings ranking. Supports `--sort`.

**Fields**: date, code, name, close, pct, hold_shares, hold_market_cap, hold_float_ratio, hold_total_ratio, add_shares, add_market_cap, add_pct

### `efi-cli stock connect hold <code> [--start YYYY-MM-DD] [--end YYYY-MM-DD]`
Individual stock holding detail.

**Fields**: date, hold_shares, hold_market_cap, hold_float_ratio, hold_total_ratio, add_shares, add_market_cap, add_pct, close, pct

### `efi-cli stock connect ah`
A+H dual-listed stocks comparison. Supports `--sort`.

**Fields**: code, name, price, pct

---

## Hot Rank Commands (人气榜)

### `efi-cli stock hot-rank list [--limit N]`
Current popularity ranking. Supports `--sort`.

**Fields**: rank, code, name, price, pct

### `efi-cli stock hot-rank history <code>`
Historical popularity trend.

**Fields**: date, rank, code, new_fans, loyal_fans

### `efi-cli stock hot-rank keyword <code>`
Hot concept keywords for a stock.

**Fields**: code, concept_name, concept_code, hot_value

---

## Comment Commands (千股千评)

### `efi-cli stock comment list [--limit N]`
Market-wide commentary scores. Supports `--sort`.

**Fields**: code, name, price, pct, turnover, pe, score, rank, attention, main_cost, institution_participation, rising, date

### `efi-cli stock comment institution <code>`
Institutional participation trend for a stock.

**Fields**: date, institution_participation

---

## Margin Commands (融资融券)

### `efi-cli stock margin account [--limit N]`
Market-wide margin trading statistics.

**Fields**: date, fin_balance, loan_balance, fin_buy_amt, loan_sell_amt, investor_num, total_guarantee, avg_guarantee_ratio, index_close, index_pct

---

## Restricted Share Commands (限售解禁)

### `efi-cli stock restricted summary [--start YYYY-MM-DD] [--end YYYY-MM-DD] [--limit N]`
Release calendar summary.

**Fields**: date, stock_count, free_shares, market_cap, index_close, index_pct

### `efi-cli stock restricted detail [--start YYYY-MM-DD] [--end YYYY-MM-DD] [--limit N]`
Detailed release list.

**Fields**: code, name, free_date, free_shares, actual_free_shares, market_cap, free_ratio, total_ratio

### `efi-cli stock restricted queue <code>`
Upcoming release batches for a stock.

**Fields**: same as restricted detail

### `efi-cli stock restricted holders <code> [--date YYYY-MM-DD]`
Release holder details.

**Fields**: holder_name, add_shares, actual_shares, add_market_cap, lock_months, residual_shares, free_type, progress

---

## Fund Commands

### `efi-cli fund quote <code1> [code2] ...`
Realtime fund NAV estimate.

**Fields**: code, name, nav, nav_date, estimate_time, estimate_pct

### `efi-cli fund info <code>`
Fund base info.

**Fields**: code, name, company, desc, estab_date, nav, nav_date, pct

### `efi-cli fund history <code> [--limit N]`
Historical NAV data.

**Fields**: code, date, nav, acc_nav, pct

### `efi-cli fund position <code> [--date YYYY-MM-DD]`
Fund stock holdings.

**Fields**: fund_code, stock_code, stock_name, pct (holding %), change (vs last period), action, sector, date

### `efi-cli fund period <code>`
Performance across time periods.

**Fields**: code, period (1w/1m/3m/6m/1y/2y/3y/5y/ytd/all), pct, avg (category average), rank, total

### `efi-cli fund asset <code> [--date YYYY-MM-DD]`
Asset allocation breakdown.

**Fields**: code, stock_pct, bond_pct, cash_pct, other_pct, total_size (亿)

### `efi-cli fund rank [--type all|gp|hh|zq] [--sort 1w|1m|3m|6m|1y|2y|3y|5y] [--limit N]`
Fund ranking by performance.

**Fields**: code, name, date, nav, acc_nav, day, week, month, m3, m6, y1, y2, y3, total, inception_date

### `efi-cli fund screen [--type all|gp|hh|zq] [--sort ...] [--size-min N] [--size-max N] [--1y-min N] [--3y-min N] [--5y-min N] [--limit N]`
Screen funds by type, size, and return thresholds.

**Parameters**:
- `--type`: all (default), gp (stock), hh (hybrid), zq (bond)
- `--sort`: 1w, 1m, 3m, 6m, 1y, 2y, 3y (default), 5y
- `--size-min/--size-max`: Fund size range (亿元)
- `--1y-min/--3y-min/--5y-min`: Min return (%) for each period

---

## Index Commands

### `efi-cli index quote <name_or_secid1> [name_or_secid2] ...`
Realtime index data. Supports Chinese names.

**Supported names**: 上证指数, 深证成指, 创业板指, 沪深300, 中证500, 中证1000, 上证50, 科创50, 恒生指数

### `efi-cli index history <name_or_secid> [--klt N] [--beg YYYYMMDD] [--end YYYYMMDD]`
Index K-line history.

---

## Bond Commands

### `efi-cli bond info <code>`
Bond base info (convertible bonds).

**Fields**: code, name, stock_code, stock_name, rating, issue_size, sub_date, list_date, expire_date, rate_desc, term

### `efi-cli bond realtime`
Realtime bond quotes.

### `efi-cli bond history <secid> [--beg YYYYMMDD] [--end YYYYMMDD]`
Bond K-line history. secid format: `market.code` (e.g., `0.127021`).

---

## Futures Commands

### `efi-cli futures realtime`
Realtime futures quotes.

### `efi-cli futures history <secid> [--beg YYYYMMDD] [--end YYYYMMDD] [--klt N]`
Futures K-line history. secid format: `exchange.code` (e.g., `113.AU` gold, `115.ZCM` thermal coal).

### `efi-cli futures deal <secid> [--max N]`
Futures tick data.

---

## Output Control Flags (all commands)

- `--format csv|json|table` — Output format (default: csv)
- `--fields f1,f2,...|all` — Select specific fields or all fields
- `--limit N` — Max rows to output
- `--no-header` — Omit header row (csv/table)
- `--pretty` — Pretty-print JSON
- `--compact` — Compact JSON
- `--schema` — Print command schema (fields, types, descriptions)
- `--list-fields` — List available field names
- `--raw` — Print raw upstream API response

## Sort Flags (selected commands)

- `--sort <field>` — Sort by field name
- `--desc` — Sort descending (default when --sort is used)
- `--asc` — Sort ascending

## Market Filter Values (for `efi-cli stock realtime <market>`)

沪深A股, 沪深京A股, 上证A股, 深证A股, 创业板, 科创板, 北证A股, 新股, 港股, 美股, ETF, 行业板块, 概念板块, 地域板块

## Market Numbers

0=深A, 1=沪A, 116/128=港股, 105/106/107=美股, 113=上期所, 114=大商所, 115=郑商所, 90=板块
