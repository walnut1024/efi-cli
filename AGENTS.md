# EFI Agent Guide

`efi-cli` queries EastMoney APIs for Chinese/US/HK stock, fund, bond, futures, and index data.

## Boundaries

- Read-only. Cannot trade or modify anything.
- All data from EastMoney public APIs only.

## Usage

```bash
make build
dist/efi-cli stock quote 600519 --format json
dist/efi-cli stock history 600519 --beg 20250101 --end 20251231
dist/efi-cli fund position 110011 --date 2025-03-31 --format json
dist/efi-cli stock screen --pe-min 5 --pe-max 20 --roe-min 15 --format json
dist/efi-cli index quote 沪深300 --format json
```

## Output

- `--format csv` (default) / `json` / `table`
- `--fields f1,f2,...` — select columns
- `--limit N` — max rows (default 50)
- `--no-header` — omit header

## Key Formats

- **secid**: `{market}.{code}` for stocks (e.g. `1.600519`), `{exchange}.{code}` for bonds/futures
- **Dates**: `YYYYMMDD` for history `--beg`/`--end`, `YYYY-MM-DD` for `--date`/`--start`/`--end`
- **K-line periods** (`--klt`): `1/5/15/30/60` min, `101`=day, `102`=week, `103`=month

## Error Handling

- Exit code 1 = failure (bad code, network error). Error to stderr.
- Empty output = no data matched (valid result).
- Client retries 3x with 1s delay on network errors.
