# Market Sentiment Reference

Use this file for sector rotation, limit-up/down pools, hot concepts, northbound flow, margin data, or unlock pressure.

## Sector and Theme

```bash
efi-cli stock sector list --type concept --sort pct --desc --limit 20
efi-cli stock sector list --type industry --sort pct --desc --limit 20
efi-cli stock sector members <code> --sort pct --desc --limit 20
efi-cli stock sector history <code> --beg YYYYMMDD
```

Look for:

- Breadth: are multiple names in the theme moving together?
- Leadership quality: is the sector led by liquid core names or thin speculative names?
- Persistence: is the move extending over days/weeks or only intraday?

## Pool and Short-Term Mood

```bash
efi-cli stock pool zt --sort seal_amount --desc --limit 20
efi-cli stock pool zt-prev --limit 20
efi-cli stock pool dt --limit 20
efi-cli stock pool zb --limit 20
efi-cli stock pool strong --limit 20
efi-cli stock hot-rank list --limit 30
```

Interpretation:

- High `seal_amount` and fewer broken boards imply stronger speculative demand
- Many broken limits imply weaker risk appetite
- Popularity should confirm price action, not replace it

## Northbound and Capital Mood

```bash
efi-cli stock connect history --type north --limit 10
efi-cli stock connect realtime
efi-cli stock connect hold-rank --type north --limit 20
efi-cli stock connect hold <code> --start YYYY-MM-DD
```

Interpretation:

- Multi-day `net_deal_amt` matters more than one strong day
- Holding-rank changes help identify accumulation, but watch whether price confirms
- Intraday flow is useful for tape color, not long-term conviction

## Margin and Unlock Pressure

```bash
efi-cli stock margin account --limit 20
efi-cli stock restricted summary --start YYYY-MM-DD --end YYYY-MM-DD --limit 20
efi-cli stock restricted detail --start YYYY-MM-DD --end YYYY-MM-DD --limit 20
efi-cli stock restricted queue <code>
```

Interpretation:

- Rising financing balance suggests stronger speculative leverage
- Falling guarantee ratio can indicate fragility
- Large upcoming unlocks can create supply pressure, especially when the stock already lacks flow support

## Good Output Shape for Sentiment Questions

1. Market mood / theme summary
2. Evidence from breadth, pools, or flow
3. What is confirming the move
4. What would signal weakening
