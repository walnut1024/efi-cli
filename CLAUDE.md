# EFI-CLI

Go CLI for querying Chinese/US/HK stock & fund data from EastMoney APIs. Output: CSV (default), JSON, or table.

## Tech Stack

- Go 1.26 + cobra CLI framework
- Module: `github.com/walnut1024/efi-cli`
- Data type: `model.Record = map[string]interface{}`, all domain funcs return `([]model.Record, error)`
- Output pipeline: `output.Format(records, format, fields, limit, noHeader)` -> stdout

## Conventions

- `push2.eastmoney.com` URLs must use `https://` (not `http://`)
- Field mappings in `internal/model/fields.go`, market filters in `internal/model/market.go`
- secid format: `{market}.{code}` (prefix `6/5/7/8`->`1.`, `0/2/3/4/9`->`0.`), bond/futures exchange-prefixed (e.g. `113.ag1024`)
- Code lookup: `client.ResolveQuoteID()` — rule-based first, then search API with structured error handling
- Build: `make build` (current platform), `make build-all` (cross-compile to dist/)
- Python: use `uv run` / `.venv/bin/python` — never `python3` or `pip install` directly

## Commands

All commands accept: `--format csv|json|table`, `--fields f1,f2,...|all`, `--limit N`, `--no-header`

Selected commands also support: `--pretty`, `--compact`, `--schema`, `--raw`, `--list-fields`

### top-level
| Command | Description |
|---------|-------------|
| `efi-cli search <keyword>` | Search securities and return candidate list |

### stock
| Command | Description |
|---------|-------------|
| `efi-cli stock quote <codes>` | Realtime quotes (supports `--sort`) |
| `efi-cli stock history <code>` | K-line history (`--klt`, `--fqt`, `--beg`, `--end`, `--indicators`, `--stats`, `--summary`) |
| `efi-cli stock realtime [market]` | Market-wide realtime (default: 沪深京A股, supports `--sort`) |
| `efi-cli stock compare <left> <right>` | Compare two stocks or indices (`--beg`, `--end`, `--metric`, `--align-date`) |
| `efi-cli stock bill <code>` | Historical capital flow |
| `efi-cli stock todaybill <code>` | Intraday capital flow |
| `efi-cli stock deal <code>` | Deal details (`--max`) |
| `efi-cli stock info <code>` | Base info (PE/PB/ROE/margins) |
| `efi-cli stock billboard` | Dragon-tiger list (`--start`, `--end`) |
| `efi-cli stock board <code>` | Boards a stock belongs to |
| `efi-cli stock members <index>` | Index constituent stocks |
| `efi-cli stock ipo` | IPO review status |
| `efi-cli stock holders` | Shareholder changes (`--date`) |
| `efi-cli stock performance` | Quarterly performance (`--date`) |
| `efi-cli stock screen` | Screen by PE/PB/ROE/cap/pct |
| `efi-cli stock finance <code>` | Financial statements (`--type income|balance|cashflow`) |
| `efi-cli stock dividend <code>` | Dividend history |
| `efi-cli stock pool zt/dt/zb/strong/sub-new/zt-prev` | Limit-up/down pools (`--date`, supports `--sort`) |
| `efi-cli stock sector list` | Sector list with quotes (`--type industry|concept`, supports `--sort`) |
| `efi-cli stock sector members <code>` | Sector constituent stocks (supports `--sort`) |
| `efi-cli stock sector history <code>` | Sector K-line history (`--klt`, `--fqt`, `--beg`, `--end`) |
| `efi-cli stock sector quote <code>` | Sector realtime quote |
| `efi-cli stock connect summary` | Stock Connect quota summary |
| `efi-cli stock connect history` | Stock Connect fund flow (`--type north|sh|sz|south`) |
| `efi-cli stock connect realtime` | Stock Connect minute-level flow |
| `efi-cli stock connect hold-rank` | Stock Connect holdings ranking (`--type`, `--indicator`, `--date`, supports `--sort`) |
| `efi-cli stock connect hold <code>` | Individual stock holding detail (`--start`, `--end`) |
| `efi-cli stock connect ah` | A+H dual-listed comparison (supports `--sort`) |
| `efi-cli stock hot-rank list` | Stock popularity ranking (supports `--sort`) |
| `efi-cli stock hot-rank history <code>` | Historical popularity trend |
| `efi-cli stock hot-rank keyword <code>` | Hot keywords for a stock |
| `efi-cli stock comment list` | Market-wide commentary scores (supports `--sort`) |
| `efi-cli stock comment institution <code>` | Institutional participation |
| `efi-cli stock margin account` | Margin trading statistics |
| `efi-cli stock restricted summary` | Restricted share release summary (`--start`, `--end`) |
| `efi-cli stock restricted detail` | Release detail (`--start`, `--end`) |
| `efi-cli stock restricted queue <code>` | Upcoming release batches |
| `efi-cli stock restricted holders <code>` | Release holder details (`--date`) |

### fund
| Command | Description |
|---------|-------------|
| `efi-cli fund history <code>` | NAV history |
| `efi-cli fund quote <codes>` | Realtime estimate |
| `efi-cli fund info <code>` | Basic info |
| `efi-cli fund position <code>` | Stock holdings (`--date`) |
| `efi-cli fund period <code>` | Period returns |
| `efi-cli fund asset <code>` | Asset allocation (`--date`) |
| `efi-cli fund rank` | Ranking (`--type all|gp|hh|zq`, `--sort 1w|1m|3m|6m|1y|2y|3y|5y`) |
| `efi-cli fund screen` | Screen funds (`--size-min/max`, `--1y-min`, `--3y-min`, `--5y-min`) |

### bond / futures / index
| Command | Description |
|---------|-------------|
| `efi-cli bond info/history/realtime` | Bond data |
| `efi-cli futures realtime/history/deal` | Futures data |
| `efi-cli index quote/history` | Index data (supports Chinese names: 上证指数, 沪深300) |
