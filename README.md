<div align="center">

# efi-cli

**A-Share / HK / US market data at your fingertips**

A blazing-fast Go CLI for querying Chinese, Hong Kong, and US stock, fund, bond,
futures & index data from public market APIs.

[![Go Report Card](https://goreportcard.com/badge/github.com/walnut1024/efi-cli)](https://goreportcard.com/report/github.com/walnut1024/efi-cli)
[![Release](https://img.shields.io/github/v/release/walnut1024/efi-cli?include_prereleases)](https://github.com/walnut1024/efi-cli/releases)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

[English](#) · [中文文档](docs/efi-cli-usage.md)

</div>

---

## Features

- **40+ commands** — quotes, K-line, capital flow, sectors, screening, margin, northbound, ...
- **Multi-asset** — stocks, funds, bonds, futures, indices
- **Flexible output** — CSV (default), JSON, or pretty table
- **AI-friendly** — structured JSON, `--schema`, `--list-fields`, `--raw` for introspection
- **Zero dependencies** — single static binary, just download and run
- **Cross-platform** — macOS (ARM/Intel), Linux, Windows

## Install

### Binary (recommended)

```bash
# macOS (Apple Silicon)
curl -sL https://github.com/walnut1024/efi-cli/releases/latest/download/efi-cli-darwin-arm64 \
  -o /usr/local/bin/efi-cli && chmod +x /usr/local/bin/efi-cli

# macOS (Intel)
curl -sL https://github.com/walnut1024/efi-cli/releases/latest/download/efi-cli-darwin-amd64 \
  -o /usr/local/bin/efi-cli && chmod +x /usr/local/bin/efi-cli

# Linux (x64)
curl -sL https://github.com/walnut1024/efi-cli/releases/latest/download/efi-cli-linux-amd64 \
  -o /usr/local/bin/efi-cli && chmod +x /usr/local/bin/efi-cli
```

<details>
<summary>Windows</summary>

Download [`efi-cli-windows-amd64.exe`](https://github.com/walnut1024/efi-cli/releases/latest) from the Releases page and add it to your `PATH`.

</details>

### From source

```bash
git clone https://github.com/walnut1024/efi-cli.git
cd efi-cli && make build
# binary → dist/efi-cli
```

## 30-Second Tour

```bash
# 🔍 Search a stock
efi-cli search 平安

# 📈 Realtime quotes (multiple at once)
efi-cli stock quote 600519 000858 000001 --format table

# 📊 K-line with MA, MACD, RSI
efi-cli stock history 600519 --beg 20240101 --end 20241231 \
  --indicators ma:5,10,20,macd,rsi:14

# 💰 Capital flow
efi-cli stock bill 600519 --limit 20

# 🌐 Northbound money
efi-cli stock connect history --type north --limit 10 --format table

# 🏭 Hot sectors
efi-cli stock sector list --type concept --sort pct --desc --limit 20

# 🎯 Screen: low PE + high ROE
efi-cli stock screen --pe-min 5 --pe-max 20 --roe-min 15

# 📉 Fund ranking
efi-cli fund rank --type gp --sort 3y --limit 10 --format table

# 📰 Index (Chinese names work)
efi-cli index quote 上证指数 沪深300 创业板指

# 🆚 Compare stock vs index
efi-cli stock compare 600519 沪深300 --beg 20240101 --metric total_return
```

## Command Reference

> All commands accept `--format csv|json|table`, `--fields`, `--limit`, `--no-header`

### Search

| Command | Description |
|---------|-------------|
| `search <keyword>` | Search by code, name, or pinyin |

### Stock

| Command | Description |
|---------|-------------|
| `stock quote <codes>` | Realtime quotes (`--sort`) |
| `stock history <code>` | K-line history (`--klt`, `--fqt`, `--beg`, `--end`, `--indicators`, `--stats`, `--summary`) |
| `stock realtime [market]` | Market-wide quotes (`--sort`) |
| `stock compare <L> <R>` | Compare two symbols (`--beg`, `--end`, `--metric`, `--align-date`) |
| `stock info <code>` | Fundamentals: PE / PB / ROE / margins |
| `stock bill <code>` | Historical daily capital flow |
| `stock todaybill <code>` | Intraday minute-level flow |
| `stock deal <code>` | Tick-by-tick deals (`--max`) |
| `stock board <code>` | Sector membership |
| `stock members <index>` | Index constituents |
| `stock billboard` | Dragon-Tiger list (`--start`, `--end`) |
| `stock ipo` | IPO review status |
| `stock holders` | Shareholder changes (`--date`) |
| `stock performance` | Quarterly results (`--date`) |
| `stock screen` | Screen by PE / PB / ROE / cap / pct |
| `stock finance <code>` | Financial statements (`--type`) |
| `stock dividend <code>` | Dividend history |

<details>
<summary>Pools, Sectors, Connect, Hot Rank, Comment, Margin, Restricted</summary>

#### Pools

| Command | Description |
|---------|-------------|
| `stock pool zt` | Limit-up stocks (`--date`, `--sort`) |
| `stock pool zt-prev` | Previous limit-up today (`--date`, `--sort`) |
| `stock pool dt` | Limit-down (`--date`, `--sort`) |
| `stock pool zb` | Broken limit (`--date`, `--sort`) |
| `stock pool strong` | Strong momentum (`--date`, `--sort`) |
| `stock pool sub-new` | Sub-new stocks (`--date`, `--sort`) |

#### Sectors

| Command | Description |
|---------|-------------|
| `stock sector list` | Industry / concept sectors (`--type`, `--sort`) |
| `stock sector members <code>` | Sector constituents (`--sort`) |
| `stock sector history <code>` | Sector K-line (`--klt`, `--beg`, `--end`) |
| `stock sector quote <code>` | Sector realtime quote |

#### Stock Connect (Northbound / Southbound)

| Command | Description |
|---------|-------------|
| `stock connect summary` | Quota overview |
| `stock connect history` | Historical fund flow (`--type north\|sh\|sz\|south`) |
| `stock connect realtime` | Minute-level flow |
| `stock connect hold-rank` | Holdings ranking (`--type`, `--indicator`, `--date`, `--sort`) |
| `stock connect hold <code>` | Individual stock holding detail (`--start`, `--end`) |
| `stock connect ah` | A+H comparison (`--sort`) |

#### Hot Rank & Commentary

| Command | Description |
|---------|-------------|
| `stock hot-rank list` | Popularity ranking (`--sort`) |
| `stock hot-rank history <code>` | Popularity trend |
| `stock hot-rank keyword <code>` | Hot concept keywords |
| `stock comment list` | Market commentary scores (`--sort`) |
| `stock comment institution <code>` | Institutional participation |

#### Margin & Restricted Shares

| Command | Description |
|---------|-------------|
| `stock margin account` | Market-wide margin statistics |
| `stock restricted summary` | Unlock calendar (`--start`, `--end`) |
| `stock restricted detail` | Unlock detail list (`--start`, `--end`) |
| `stock restricted queue <code>` | Upcoming unlocks |
| `stock restricted holders <code>` | Unlock holder details (`--date`) |

</details>

### Fund

| Command | Description |
|---------|-------------|
| `fund quote <codes>` | Realtime NAV estimate |
| `fund history <code>` | NAV history |
| `fund info <code>` | Fund basics |
| `fund position <code>` | Stock holdings (`--date`) |
| `fund period <code>` | Period returns |
| `fund asset <code>` | Asset allocation (`--date`) |
| `fund rank` | Ranking (`--type`, `--sort`) |
| `fund screen` | Screen funds (`--size-min/max`, `--1y-min`, `--3y-min`, `--5y-min`) |

### Bond / Futures / Index

| Command | Description |
|---------|-------------|
| `bond info\|history\|realtime` | Bond data (convertible bonds) |
| `futures realtime\|history\|deal` | Futures data |
| `index quote\|history` | Index data (Chinese names: 上证指数, 沪深300 ...) |

## Output Control

| Flag | Description |
|------|-------------|
| `--format csv\|json\|table` | Output format (default: csv) |
| `--fields f1,f2,...` | Select columns (`all` for everything) |
| `--limit N` | Max rows (default 50) |
| `--no-header` | Omit header row |
| `--sort <field>` | Sort by field (`--desc` default / `--asc`) |
| `--pretty` | Pretty-print JSON |
| `--schema` | Print command field schema |
| `--list-fields` | List all available fields |
| `--raw` | Dump raw upstream API response |

## Build

```bash
make build        # Current platform  → dist/efi-cli
make build-all    # Cross-compile     → dist/
make test         # Run tests
make clean        # Remove dist/
```

## License

[MIT](LICENSE)
