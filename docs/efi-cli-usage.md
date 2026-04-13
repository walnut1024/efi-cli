# EFI CLI 使用文档

## 简介

`efi-cli` 是一个面向终端和 AI agent 场景的中国金融市场数据 CLI 工具。

支持的数据类型：

1. A 股（沪深京）、港股通、A+H 股
2. 公募基金（股票型 / 混合型 / 债券型）
3. 债券（可转债等）
4. 期货
5. 主要指数（上证、深证、创业板、沪深 300 等）

典型使用场景：

1. 终端内快速查看行情、K 线、资金流向、板块数据。
2. 在 AI agent 或脚本中以 JSON 格式消费数据。
3. 通过 `--schema`、`--list-fields`、`--raw` 做能力探测和调试。

## 安装与构建

```bash
make build
```

构建完成后：

```bash
./dist/efi-cli --help
```

开发模式：

```bash
go run ./cmd/efi-cli --help
```

## 命令结构总览

```text
efi-cli
├── search                          搜索证券候选
├── stock                           股票市场数据
│   ├── quote                       实时行情
│   ├── history                     历史 K 线
│   ├── realtime                    全市场实时行情
│   ├── compare                     两标的对比
│   ├── info                        基本信息
│   ├── bill                        历史资金流向
│   ├── todaybill                   分时资金流向
│   ├── deal                        成交明细
│   ├── board                       所属板块
│   ├── members                     指数成分股
│   ├── billboard                   龙虎榜
│   ├── ipo                         IPO 审核
│   ├── holders                     股东人数变化
│   ├── performance                 业绩报表
│   ├── screen                      条件选股
│   ├── finance                     财务报表
│   ├── dividend                    分红送转
│   ├── pool                        涨跌停 / 强势 / 次新池
│   │   ├── zt                      涨停池
│   │   ├── zt-prev                 昨日涨停
│   │   ├── dt                      跌停池
│   │   ├── zb                      炸板池
│   │   ├── strong                  强势股池
│   │   └── sub-new                 次新股池
│   ├── sector                      板块深数据
│   │   ├── list                    板块列表
│   │   ├── members                 板块成分股
│   │   ├── history                 板块 K 线
│   │   └── quote                   板块实时行情
│   ├── connect                     沪深港通
│   │   ├── summary                 通道路径汇总
│   │   ├── history                 历史资金流向
│   │   ├── realtime                分时资金流向
│   │   ├── hold-rank               持股排行
│   │   ├── hold                    个股持股明细
│   │   └── ah                      A+H 对比
│   ├── hot-rank                    人气榜
│   │   ├── list                    当前人气排名
│   │   ├── history                 历史人气趋势
│   │   └── keyword                 热门概念关键词
│   ├── comment                     千股千评
│   │   ├── list                    市场评分排行
│   │   └── institution             个股机构参与度
│   ├── margin                      融资融券
│   │   └── account                 全市场融资融券统计
│   └── restricted                  限售解禁
│       ├── summary                 解禁日历汇总
│       ├── detail                  解禁明细列表
│       ├── queue                   个股未来解禁排期
│       └── holders                 解禁股东明细
├── fund                            基金数据
│   ├── history                     净值历史
│   ├── quote                       实时估值
│   ├── info                        基本信息
│   ├── position                    持仓明细
│   ├── period                      阶段收益
│   ├── asset                       资产配置
│   ├── rank                        基金排名
│   └── screen                      基金筛选
├── bond                            债券数据
│   ├── info                        债券信息
│   ├── history                     K 线历史
│   └── realtime                    实时行情
├── futures                         期货数据
│   ├── realtime                    实时行情
│   ├── history                     K 线历史
│   └── deal                        成交明细
└── index                           指数数据
    ├── quote                       实时指数
    └── history                     指数 K 线
```

## 通用输出参数

所有命令都支持以下通用参数：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--format` | `csv` | 输出格式：`csv`、`json`、`table` |
| `--fields` | （各命令默认字段集） | 指定输出字段，逗号分隔；`all` 输出全部字段 |
| `--limit` | `50` | 最大返回行数 |
| `--no-header` | `false` | 省略表头（csv / table 格式） |

部分命令额外支持的高级参数：

| 参数 | 说明 |
|------|------|
| `--pretty` | 缩进 JSON 输出，适合人工阅读 |
| `--compact` | 紧凑 JSON 输出 |
| `--schema` | 输出命令的 schema 定义（字段、类型、描述），不返回业务数据 |
| `--list-fields` | 输出当前命令支持的全部字段名 |
| `--raw` | 输出上游原始接口 JSON 响应，适合调试 |

支持排序的命令额外支持：

| 参数 | 说明 |
|------|------|
| `--sort <field>` | 按指定字段排序 |
| `--desc` | 降序（默认） |
| `--asc` | 升序 |

排序规则：

1. 数值字段按数值排序，空值排最后。
2. `--asc` 和 `--desc` 不能同时使用。
3. 必须搭配 `--sort` 使用。

---

## 搜索命令

### search — 搜索证券候选

```bash
efi-cli search <keyword>
```

输入股票简称、代码、拼音首字母、指数名称等关键字，返回候选列表。

```bash
# 基本搜索
efi-cli search 平安
efi-cli search 600519
efi-cli search pingan

# 指定格式和字段
efi-cli search 沪深300 --format json --pretty
efi-cli search 159915 --fields code,name,market,quote_id

# 能力探测
efi-cli search pingan --schema --pretty
efi-cli search pingan --list-fields
```

返回字段：`code`、`name`、`market`、`quote_id`、`secid` 等。

---

## 股票命令

### stock quote — 实时行情

查询一只或多只股票的实时行情。支持 `--sort`。

```bash
# 单只
efi-cli stock quote 600519

# 多只
efi-cli stock quote 600519 000858 000001

# 指定字段和格式
efi-cli stock quote 600519 --format json --pretty
efi-cli stock quote 600519 000858 --fields code,name,price,pct,amount

# 排序
efi-cli stock quote 600519 000858 000001 --sort pct --desc
```

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | string | 证券代码 |
| `name` | string | 证券名称 |
| `price` | number | 最新价 |
| `pct` | number | 涨跌幅 (%) |
| `chg` | number | 涨跌额 |
| `open` | number | 开盘价 |
| `high` | number | 最高价 |
| `low` | number | 最低价 |
| `pre_close` | number | 昨收价 |
| `vol` | number | 成交量（手） |
| `amount` | number | 成交额 |
| `turnover` | number | 换手率 (%) |
| `vol_ratio` | number | 量比 |
| `pe` | number | 市盈率（动态） |
| `pb` | number | 市净率 |
| `roe` | number | ROE |
| `mkt_cap` | number | 总市值 |
| `float_cap` | number | 流通市值 |
| `industry` | string | 所属行业 |
| `board_name` | string | 所属板块 |
| `latest_trade_date` | string | 最近交易日期 |
| `secid` | string | secid |

### stock realtime — 全市场实时行情

按市场板块获取全市场实时行情。支持 `--sort`。

```bash
efi-cli stock realtime                        # 默认：沪深京 A 股
efi-cli stock realtime 创业板                  # 创业板
efi-cli stock realtime 沪深京A股 --sort amount --desc --limit 20
efi-cli stock realtime 科创板 --sort turnover --asc --format table
```

可用市场名：`沪深京A股`、`沪深A股`、`创业板`、`科创板` 等。

字段与 `stock quote` 相同。

### stock history — 历史 K 线

```bash
efi-cli stock history <code> [--beg YYYYMMDD] [--end YYYYMMDD] [--klt 101] [--fqt 1] \
  [--indicators ...] [--stats ...] [--summary]
```

参数：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--klt` | `101` | K 线周期：`1`/`5`/`15`/`30`/`60` 分钟，`101` 日线，`102` 周线，`103` 月线 |
| `--fqt` | `1` | 复权方式：`0` 不复权，`1` 前复权，`2` 后复权 |
| `--beg` | `19000101` | 起始日期 YYYYMMDD |
| `--end` | `20500101` | 结束日期 YYYYMMDD |
| `--indicators` | 无 | 附加技术指标，见下表 |
| `--stats` | 无 | 区间统计项，见下表 |
| `--summary` | `false` | 返回单行摘要而不是逐行数据 |

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 证券名称 |
| `code` | string | 证券代码 |
| `date` | string | 日期 |
| `open` | number | 开盘价 |
| `close` | number | 收盘价 |
| `high` | number | 最高价 |
| `low` | number | 最低价 |
| `vol` | number | 成交量 |
| `amount` | number | 成交额 |
| `pct` | number | 涨跌幅 (%) |
| `chg` | number | 涨跌额 |
| `turnover` | number | 换手率 (%) |
| `amp` | number | 振幅 (%) |

技术指标（`--indicators`）：

| 指标 | 语法 | 生成字段 |
|------|------|----------|
| MA | `ma:5,10,20` | `ma5`, `ma10`, `ma20` |
| EMA | `ema:12,26` | `ema12`, `ema26` |
| MACD | `macd` | `dif`, `dea`, `macd` |
| RSI | `rsi:14` | `rsi14` |
| BOLL | `boll:20,2` | `boll_mid`, `boll_up`, `boll_low` |

区间统计（`--stats`）：

| 统计项 | 说明 |
|--------|------|
| `total_return` | 区间总收益率 (%) |
| `cumulative_pct` | 累计收益率 (%) |
| `amplitude_avg` | 平均振幅 |
| `high` | 区间最高价 |
| `low` | 区间最低价 |
| `max_drawdown` | 最大回撤 (%) |
| `start_close` | 起始收盘价 |
| `end_close` | 结束收盘价 |
| `start_date` | 起始日期 |
| `end_date` | 结束日期 |
| `bars` | K 线根数 |

示例：

```bash
# 最近日线
efi-cli stock history 600519

# 指定日期范围
efi-cli stock history 600519 --beg 20240101 --end 20241231

# 周线 + 前复权
efi-cli stock history 600519 --klt 102 --fqt 1

# 附加技术指标
efi-cli stock history 600519 --indicators ma:5,10,20,macd,rsi:14

# 区间统计摘要
efi-cli stock history 600519 --beg 20240101 --end 20241231 \
  --stats total_return,high,low,max_drawdown --summary --format json

# 组合：指标 + 统计
efi-cli stock history 600519 --indicators ma:20 --stats cumulative_pct
```

### stock compare — 两标的对比

对比两只股票或指数在指定区间的表现。

```bash
efi-cli stock compare <left> <right> [--beg] [--end] [--metric] [--align-date]
```

参数：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--klt` | `101` | K 线周期：`101` 日线，`102` 周线，`103` 月线 |
| `--beg` | `19000101` | 起始日期 |
| `--end` | `20500101` | 结束日期 |
| `--metric` | 无 | 核心比较指标：`total_return`, `max_drawdown`, `high`, `low`, `amplitude_avg` 等 |
| `--align-date` | `false` | 按共同日期对齐输出双边时间序列 |

支持股票和内置指数名混合比较（如 `沪深300`、`上证指数`）。支持 `--sort`。

```bash
# 基本对比
efi-cli stock compare 600519 000858

# 指定区间
efi-cli stock compare 600519 沪深300 --beg 20240101 --end 20241231

# 指定核心指标
efi-cli stock compare 600519 沪深300 --metric total_return

# 按日期对齐输出
efi-cli stock compare 600519 沪深300 --align-date --format json --pretty
```

### stock info — 基本信息

获取股票基本面数据（PE/PB/ROE/利润率等）。

```bash
efi-cli stock info 600519
```

输出字段：`code`, `name`, `industry`, `board_name`, `price`, `pct`, `pe`, `pb`, `roe`, `net_profit`, `gross_margin`, `net_margin`, `mkt_cap`, `float_cap`。

### stock bill — 历史资金流向

获取个股历史每日的主力/散户资金流向。

```bash
efi-cli stock bill 600519
```

输出字段：`date`, `main_net`, `main_pct`, `small_net`, `small_pct`, `med_net`, `med_pct`, `big_net`, `big_pct`, `huge_net`, `huge_pct`。

### stock todaybill — 分时资金流向

获取个股当日分钟级资金流向。

```bash
efi-cli stock todaybill 600519
```

输出字段：`time`, `main_net`, `small_net`, `med_net`, `big_net`, `huge_net`。

### stock deal — 成交明细

获取最近交易日的逐笔成交。

```bash
efi-cli stock deal 600519 --max 500
```

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--max` | `10000` | 最大返回笔数 |

输出字段：`code`, `time`, `price`, `vol`, `num`, `pre_close`。

### stock board — 所属板块

查看一只股票所属的所有板块。

```bash
efi-cli stock board 600519
```

### stock members — 指数成分股

查看指数的成分股列表。

```bash
efi-cli stock members 沪深300
```

支持内置指数名：`上证指数`、`深证成指`、`创业板指`、`沪深300`、`中证500`、`中证1000`、`上证50`、`科创50`、`恒生指数`。

### stock billboard — 龙虎榜

```bash
efi-cli stock billboard --start 2024-01-01 --end 2024-01-31
```

| 参数 | 说明 |
|------|------|
| `--start` | 起始日期 YYYY-MM-DD |
| `--end` | 结束日期 YYYY-MM-DD |

### stock ipo — IPO 审核状态

```bash
efi-cli stock ipo
```

### stock holders — 股东人数变化

```bash
efi-cli stock holders --date 2024-12-31
```

| 参数 | 说明 |
|------|------|
| `--date` | 报告期 YYYY-MM-DD |

### stock performance — 业绩报表

```bash
efi-cli stock performance --date 2024-12-31
```

| 参数 | 说明 |
|------|------|
| `--date` | 报告期 YYYY-MM-DD |

### stock screen — 条件选股

按 PE、PB、ROE、市值、涨跌幅等条件筛选股票。

```bash
efi-cli stock screen [--market 沪深京A股] [--pe-min] [--pe-max] [--pb-min] [--pb-max] \
  [--roe-min] [--cap-min] [--cap-max] [--pct-min] [--pct-max] [--limit 50]
```

参数：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--market` | `沪深京A股` | 市场范围：`沪深A股`、`创业板`、`科创板` 等 |
| `--pe-min` | `0` | 最小市盈率 |
| `--pe-max` | `0` | 最大市盈率（0 = 不限） |
| `--pb-min` | `0` | 最小市净率 |
| `--pb-max` | `0` | 最大市净率（0 = 不限） |
| `--roe-min` | `0` | 最小 ROE (%) |
| `--cap-min` | `0` | 最小市值（亿元） |
| `--cap-max` | `0` | 最大市值（亿元，0 = 不限） |
| `--pct-min` | `0` | 最小当日涨跌幅 (%) |
| `--pct-max` | `0` | 最大当日涨跌幅 (%) |

```bash
# 全部 A 股，默认条件
efi-cli stock screen

# 低 PE + 高 ROE
efi-cli stock screen --pe-min 5 --pe-max 20 --roe-min 15

# 大盘蓝筹
efi-cli stock screen --cap-min 500 --cap-max 5000 --pe-max 30

# 限范围输出
efi-cli stock screen --market 创业板 --limit 20 --format table
```

### stock finance — 财务报表

```bash
efi-cli stock finance <code> --type income|balance|cashflow
```

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--type` | `income` | 报表类型：`income`（利润表）、`balance`（资产负债表）、`cashflow`（现金流量表） |

### stock dividend — 分红送转

```bash
efi-cli stock dividend 600519
```

---

## 涨跌停池命令

`stock pool` 命令组提供 6 个子命令，用于查看当日各类异动股票池。所有子命令均支持 `--sort` 和 `--schema`。

### stock pool zt — 涨停池

当日涨停股票。

```bash
efi-cli stock pool zt --date 20260408
efi-cli stock pool zt --sort seal_amount --desc --limit 20
efi-cli stock pool zt --format table
```

| 参数 | 说明 |
|------|------|
| `--date` | 日期 YYYYMMDD（默认当天） |

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | string | 证券代码 |
| `name` | string | 证券名称 |
| `pct` | number | 涨跌幅 |
| `price` | number | 最新价 |
| `amount` | number | 成交额 |
| `float_cap` | number | 流通市值 |
| `total_cap` | number | 总市值 |
| `turnover` | number | 换手率 |
| `seal_amount` | number | 封单金额 |
| `first_seal_time` | string | 首次封板时间 |
| `last_seal_time` | string | 最后封板时间 |
| `break_count` | number | 炸板次数 |
| `consecutive` | number | 连板数 |
| `industry` | string | 所属行业 |

### stock pool zt-prev — 昨日涨停

前一个交易日涨停的股票今日表现。

```bash
efi-cli stock pool zt-prev --date 20260408
```

输出字段与涨停池相同（无 `seal_amount`、`first_seal_time`、`break_count`、`consecutive`）。

### stock pool dt — 跌停池

当日跌停股票。

```bash
efi-cli stock pool dt --date 20260408
```

输出字段：在通用池字段基础上增加 `pe`（市盈率）、`seal_amount`（封单金额）、`last_seal_time`、`break_count`、`consecutive`。

### stock pool zb — 炸板池

曾涨停但未能封住的股票。

```bash
efi-cli stock pool zb --date 20260408
```

### stock pool strong — 强势股池

当日表现强势的股票。

```bash
efi-cli stock pool strong --date 20260408
```

### stock pool sub-new — 次新股池

近期上市的次新股。

```bash
efi-cli stock pool sub-new --date 20260408
```

---

## 板块深数据命令

`stock sector` 命令组提供行业/概念板块的深度数据。

### stock sector list — 板块列表

获取行业或概念板块列表及实时行情。支持 `--sort`。

```bash
efi-cli stock sector list --type concept
efi-cli stock sector list --type industry --sort pct --desc --limit 30
efi-cli stock sector list --format table
```

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--type` | `concept` | 板块类型：`industry`（行业）、`concept`（概念） |

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | string | 板块代码 |
| `name` | string | 板块名称 |
| `price` | number | 板块指数 |
| `pct` | number | 涨跌幅 |
| `chg` | number | 涨跌额 |
| `total_cap` | number | 总市值 |
| `turnover` | number | 换手率 |
| `up_count` | number | 上涨家数 |
| `down_count` | number | 下跌家数 |
| `lead_stock` | string | 领涨股名称 |
| `lead_code` | string | 领涨股代码 |

### stock sector members — 板块成分股

获取指定板块的成分股列表。支持 `--sort`。

```bash
efi-cli stock sector members BK0477
efi-cli stock sector members BK0477 --sort pct --desc --limit 30
```

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | string | 证券代码 |
| `name` | string | 证券名称 |
| `price` | number | 最新价 |
| `pct` | number | 涨跌幅 |
| `chg` | number | 涨跌额 |
| `vol` | number | 成交量 |
| `amt` | number | 成交额 |
| `amplitude` | number | 振幅 |
| `high` | number | 最高价 |
| `low` | number | 最低价 |
| `open` | number | 开盘价 |
| `pre_close` | number | 昨收价 |
| `turnover` | number | 换手率 |
| `pe` | number | 市盈率 |
| `pb` | number | 市净率 |

### stock sector history — 板块 K 线

获取板块指数的 K 线历史数据。

```bash
efi-cli stock sector history BK0477 --beg 20240101 --end 20241231
efi-cli stock sector history BK0477 --klt 102 --fqt 1
```

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--klt` | `101` | 周期：`101` 日线，`102` 周线，`103` 月线 |
| `--fqt` | `1` | 复权：`0` 不复权，`1` 前复权 |
| `--beg` | `19000101` | 起始日期 YYYYMMDD |
| `--end` | `20500101` | 结束日期 YYYYMMDD |

### stock sector quote — 板块实时行情

```bash
efi-cli stock sector quote BK0477
```

返回该板块的实时行情数据（同 `sector list` 的字段）。

---

## 沪深港通命令

`stock connect` 命令组提供沪深港通（北向/南向）资金数据。

### stock connect summary — 通道路径汇总

各沪深港通渠道的额度汇总。

```bash
efi-cli stock connect summary
efi-cli stock connect summary --format table
efi-cli stock connect summary --schema --pretty
```

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `date` | string | 交易日期 |
| `mutual_type_name` | string | 通道路径名称 |
| `direction` | string | 资金方向 |
| `board_type` | string | 板块类型 |
| `index_name` | string | 关联指数 |

### stock connect history — 历史资金流向

沪深港通每日资金流入流出历史。

```bash
efi-cli stock connect history --type north --limit 30
efi-cli stock connect history --type sh --format json --pretty
efi-cli stock connect history --type south --limit 10
```

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--type` | `north` | 类型：`north`（北向合计）、`sh`（沪股通）、`sz`（深股通）、`south`（南向）、`hshk`（港股通沪港深）、`szhk`（港股通深港深） |

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `date` | string | 交易日期 |
| `fund_inflow` | number | 资金流入 |
| `net_deal_amt` | number | 净成交额 |
| `buy_amt` | number | 买入额 |
| `sell_amt` | number | 卖出额 |
| `hold_market_cap` | number | 持股市值 |
| `quota_balance` | number | 余额 |
| `lead_stock` | string | 领涨股 |
| `lead_code` | string | 领涨股代码 |

### stock connect realtime — 分时资金流向

沪深港通当日分钟级资金流向。

```bash
efi-cli stock connect realtime
efi-cli stock connect realtime --format table
```

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `time` | string | 时间 |
| `sh_connect` | number | 沪股通（万元） |
| `sz_connect` | number | 深股通（万元） |
| `north_total` | number | 北向合计（万元） |

### stock connect hold-rank — 持股排行

沪深港通持股量排名。支持 `--sort`。

```bash
efi-cli stock connect hold-rank --type north --limit 30
efi-cli stock connect hold-rank --type sh --indicator 1 --date 2026-04-07
efi-cli stock connect hold-rank --sort hold_market_cap --desc --limit 20
```

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--type` | `north` | 类型：`north`（北向）、`sh`（沪股通）、`sz`（深股通） |
| `--indicator` | 无 | 周期：`1`（1日）、`3`、`5`、`10`、`M`（月）、`Q`（季）、`Y`（年） |
| `--date` | 无 | 日期 YYYY-MM-DD |

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `date` | string | 日期 |
| `code` | string | 证券代码 |
| `name` | string | 证券名称 |
| `close` | number | 收盘价 |
| `pct` | number | 涨跌幅 |
| `hold_shares` | number | 持股数量 |
| `hold_market_cap` | number | 持股市值 |
| `hold_float_ratio` | number | 占流通股比 |
| `hold_total_ratio` | number | 占总股本比 |
| `add_shares` | number | 增持股数 |
| `add_market_cap` | number | 增持市值 |
| `add_pct` | number | 增持比例 |

### stock connect hold — 个股持股明细

查看某只股票被沪深港通持有的历史变化。

```bash
efi-cli stock connect hold 600519
efi-cli stock connect hold 600519 --start 2026-01-01 --end 2026-03-31
```

| 参数 | 说明 |
|------|------|
| `--start` | 起始日期 YYYY-MM-DD |
| `--end` | 结束日期 YYYY-MM-DD |

输出字段：`date`, `hold_shares`, `hold_market_cap`, `hold_float_ratio`, `hold_total_ratio`, `add_shares`, `add_market_cap`, `add_pct`, `close`, `pct`。

### stock connect ah — A+H 对比

A+H 上市股票的对比数据。支持 `--sort`。

```bash
efi-cli stock connect ah
efi-cli stock connect ah --sort pct --desc --format table
```

输出字段：`code`, `name`, `price`, `pct`。

---

## 人气榜命令

`stock hot-rank` 命令组提供股票人气排名数据。

### stock hot-rank list — 当前人气排名

获取全市场人气排名。支持 `--sort`。

```bash
efi-cli stock hot-rank list --limit 30
efi-cli stock hot-rank list --sort rank --asc --format table
```

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `rank` | number | 排名 |
| `code` | string | 证券代码 |
| `name` | string | 证券名称 |
| `price` | number | 最新价 |
| `pct` | number | 涨跌幅 |

### stock hot-rank history — 历史人气趋势

查看单只股票的历史人气排名变化。

```bash
efi-cli stock hot-rank history 600519
```

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `date` | string | 日期 |
| `rank` | number | 排名 |
| `code` | string | 证券代码 |
| `new_fans` | number | 新增粉丝 |
| `loyal_fans` | number | 忠诚粉丝 |

### stock hot-rank keyword — 热门概念关键词

查看与某只股票相关的热门概念。

```bash
efi-cli stock hot-rank keyword 600519
```

输出字段：`code`, `concept_name`, `concept_code`, `hot_value`。

---

## 千股千评命令

`stock comment` 命令组提供股票综合评分数据。

### stock comment list — 市场评分排行

获取全市场股票的综合评分排行。支持 `--sort`。

```bash
efi-cli stock comment list --limit 30
efi-cli stock comment list --sort score --desc --format table
efi-cli stock comment list --sort rank --asc --limit 50
```

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | string | 证券代码 |
| `name` | string | 证券名称 |
| `price` | number | 收盘价 |
| `pct` | number | 涨跌幅 |
| `turnover` | number | 换手率 |
| `pe` | number | 市盈率 |
| `score` | number | 综合评分 |
| `rank` | number | 排名 |
| `attention` | number | 关注度 |
| `main_cost` | number | 主力成本 |
| `institution_participation` | number | 机构参与度 |
| `rising` | number | 上升比例 |
| `date` | string | 日期 |

### stock comment institution — 个股机构参与度

查看单只股票的机构参与度历史趋势。

```bash
efi-cli stock comment institution 600519
```

输出字段：`date`, `institution_participation`。

---

## 融资融券命令

### stock margin account — 全市场融资融券统计

获取全市场融资融券余额等统计数据。

```bash
efi-cli stock margin account --limit 30
efi-cli stock margin account --format table
efi-cli stock margin account --schema --pretty
```

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `date` | string | 统计日期 |
| `fin_balance` | number | 融资余额（亿元） |
| `loan_balance` | number | 融券余额（亿元） |
| `fin_buy_amt` | number | 融资买入额（亿元） |
| `loan_sell_amt` | number | 融券卖出额（亿元） |
| `investor_num` | number | 参与交易投资者数 |
| `total_guarantee` | number | 担保物总价值（亿元） |
| `avg_guarantee_ratio` | number | 平均维持担保比例 |
| `index_close` | number | 上证收盘 |
| `index_pct` | number | 上证涨跌幅 |

---

## 限售解禁命令

`stock restricted` 命令组提供限售股解禁数据。

### stock restricted summary — 解禁日历汇总

按日期汇总的解禁统计。

```bash
efi-cli stock restricted summary --start 2026-01-01 --end 2026-06-30
efi-cli stock restricted summary --limit 30 --format table
```

| 参数 | 说明 |
|------|------|
| `--start` | 起始日期 YYYY-MM-DD |
| `--end` | 结束日期 YYYY-MM-DD |

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `date` | string | 解禁日期 |
| `stock_count` | number | 解禁公司数 |
| `free_shares` | number | 解禁股数（万股） |
| `market_cap` | number | 解禁市值（万元） |
| `index_close` | number | 指数收盘 |
| `index_pct` | number | 指数涨跌幅 |

### stock restricted detail — 解禁明细列表

个股级别的解禁详情。

```bash
efi-cli stock restricted detail --start 2026-04-01 --end 2026-04-30
```

| 参数 | 说明 |
|------|------|
| `--start` | 起始日期 YYYY-MM-DD |
| `--end` | 结束日期 YYYY-MM-DD |

输出字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | string | 证券代码 |
| `name` | string | 证券名称 |
| `free_date` | string | 解禁日期 |
| `free_shares` | number | 计划解禁股数 |
| `actual_free_shares` | number | 实际解禁股数 |
| `market_cap` | number | 解禁市值 |
| `free_ratio` | number | 解禁比例 |
| `total_ratio` | number | 占总股本比 |

### stock restricted queue — 个股未来解禁排期

查看某只股票未来的解禁批次。

```bash
efi-cli stock restricted queue 600519
```

输出字段同 `restricted detail`。

### stock restricted holders — 解禁股东明细

查看某次解禁涉及的股东详情。

```bash
efi-cli stock restricted holders 600519 --date 2026-04-07
```

| 参数 | 说明 |
|------|------|
| `--date` | 解禁日期 YYYY-MM-DD |

输出字段：`holder_name`, `add_shares`, `actual_shares`, `add_market_cap`, `lock_months`, `residual_shares`, `free_type`, `progress`。

---

## 指数命令

### index quote — 实时指数

```bash
efi-cli index quote 上证指数 沪深300
efi-cli index quote 1.000001 0.399001
```

支持内置指数名：

| 指数名 | secid |
|--------|-------|
| `上证指数` | `1.000001` |
| `深证成指` | `0.399001` |
| `创业板指` | `0.399006` |
| `沪深300` | `1.000300` |
| `中证500` | `1.000905` |
| `中证1000` | `0.000852` |
| `上证50` | `1.000016` |
| `科创50` | `1.000688` |
| `恒生指数` | `100.HSI` |

### index history — 指数 K 线

```bash
efi-cli index history 沪深300 --beg 20240101 --end 20241231
efi-cli index history 上证指数 --klt 102
```

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--klt` | `101` | 周期：`101` 日线，`102` 周线，`103` 月线 |
| `--beg` | `19000101` | 起始日期 YYYYMMDD |
| `--end` | `20500101` | 结束日期 YYYYMMDD |

---

## 基金命令

### fund quote — 实时估值

```bash
efi-cli fund quote 161725 000311
efi-cli fund quote 161725 --format json --pretty
```

输出字段：`code`, `name`, `nav`, `nav_date`, `estimate_time`, `estimate_pct`。

### fund history — 净值历史

```bash
efi-cli fund history 161725
```

输出字段：`code`, `date`, `nav`, `acc_nav`, `pct`。

### fund info — 基本信息

```bash
efi-cli fund info 161725
```

输出字段：`code`, `name`, `estab_date`, `pct`, `nav`, `company`, `nav_date`, `desc`。

### fund position — 持仓明细

```bash
efi-cli fund position 161725 --date 2024-12-31
```

| 参数 | 说明 |
|------|--------|
| `--date` | 报告期 YYYY-MM-DD |

输出字段：`fund_code`, `stock_code`, `stock_name`, `pct`, `change`, `action`, `sector`, `date`。

### fund period — 阶段收益

```bash
efi-cli fund period 161725
```

输出字段：`code`, `pct`, `avg`, `rank`, `total`, `period`。

`period` 值：`1w`, `1m`, `3m`, `6m`, `1y`, `2y`, `3y`, `5y`, `ytd`, `all`。

### fund asset — 资产配置

```bash
efi-cli fund asset 161725 --date 2024-12-31
```

| 参数 | 说明 |
|------|--------|
| `--date` | 报告期 YYYY-MM-DD |

输出字段：`code`, `stock_pct`, `bond_pct`, `cash_pct`, `other_pct`, `total_size`。

### fund rank — 基金排名

```bash
efi-cli fund rank --type gp --sort 3y
efi-cli fund rank --type all --sort 1y --limit 30
```

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--type` | `all` | 基金类型：`all`（全部）、`gp`（股票型）、`hh`（混合型）、`zq`（债券型） |
| `--sort` | `3y` | 排序指标：`1w`, `1m`, `3m`, `6m`, `1y`, `2y`, `3y`, `5y` |

输出字段：`code`, `name`, `date`, `nav`, `acc_nav`, `day`, `week`, `month`, `m3`, `m6`, `y1`, `y2`, `y3`, `total`, `inception_date`。

### fund screen — 基金筛选

按基金类型、规模、收益率等条件筛选。

```bash
efi-cli fund screen --type gp --size-min 10 --1y-min 5 --3y-min 10
efi-cli fund screen --sort 5y --5y-min 50 --limit 30
efi-cli fund screen --size-min 5 --size-max 50 --format table
```

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--type` | `all` | 基金类型：`all`、`gp`、`hh`、`zq` |
| `--sort` | `3y` | 排序指标：`1w`, `1m`, `3m`, `6m`, `1y`, `2y`, `3y`, `5y` |
| `--size-min` | `0` | 最小基金规模（亿元） |
| `--size-max` | `0` | 最大基金规模（亿元，0 = 不限） |
| `--1y-min` | `0` | 最小近 1 年收益率 (%) |
| `--3y-min` | `0` | 最小近 3 年收益率 (%) |
| `--5y-min` | `0` | 最小近 5 年收益率 (%) |

---

## 债券命令

```bash
efi-cli bond info 113527
efi-cli bond realtime
efi-cli bond history 113.ag1024 --beg 20240101 --end 20241231
```

说明：`bond history` 接受 `secid`，格式为 `{交易所代码}.{代码}`（如 `113.ag1024`）。

---

## 期货命令

```bash
efi-cli futures realtime
efi-cli futures history 113.ag1024 --beg 20240101 --end 20241231
efi-cli futures deal 115.ZCM --max 500
```

说明：`futures history` 和 `futures deal` 接受 `secid`。

---

## Schema 与字段探测

在 AI agent 或脚本中，建议先做能力探测：

```bash
# 查看命令 schema（字段定义 + 支持的功能）
efi-cli stock quote --schema
efi-cli stock history --schema --pretty
efi-cli stock pool zt --schema --pretty
efi-cli stock sector list --schema --pretty
efi-cli stock connect hold-rank --schema --pretty
efi-cli stock comment list --schema --pretty
efi-cli stock margin account --schema --pretty
efi-cli stock restricted detail --schema --pretty

# 仅列出字段名
efi-cli stock history --list-fields
efi-cli stock pool zt --list-fields
efi-cli stock sector members --list-fields

# 输出全部字段
efi-cli stock quote 600519 --fields all
efi-cli stock history 600519 --fields all --format json --pretty
```

## 原始响应调试

使用 `--raw` 可直接输出上游接口的原始 JSON：

```bash
efi-cli search 平安 --raw
efi-cli stock quote 600519 --raw
efi-cli stock history 600519 --raw
efi-cli stock realtime --raw
efi-cli stock pool zt --raw
```

常见用途：

1. 调试字段映射是否正确。
2. 回归对比上游接口是否变动。
3. 为新增字段采集样本数据。

## 错误输出

CLI 区分以下错误类型：

| 错误类型 | 说明 |
|----------|------|
| `代码未找到` | 输入的证券代码无匹配 |
| `匹配到多个标的` | 搜索结果有歧义 |
| `上游接口异常` | 上游接口返回错误或网络超时 |
| `上游响应解析异常` | 接口返回数据格式异常 |
| `参数错误` | 命令行参数不合法 |

在 JSON 模式下，错误以结构化对象输出：

```json
{
  "error": {
    "kind": "invalid_argument",
    "message": "参数错误",
    "input": "pct",
    "op": "sort"
  }
}
```

## 典型用法示例

### 1. 先搜索，再查询

```bash
efi-cli search 平安 --format table
efi-cli stock quote 000001 --format table
```

### 2. 查看市场成交额 Top 20

```bash
efi-cli stock realtime 沪深京A股 --sort amount --desc --limit 20 --format table
```

### 3. 技术指标分析

```bash
efi-cli stock history 600519 --beg 20240101 --end 20241231 \
  --indicators ma:5,10,20,macd,rsi:14
```

### 4. 区间摘要

```bash
efi-cli stock history 600519 --beg 20240101 --end 20241231 \
  --stats total_return,high,low,max_drawdown --summary --format json
```

### 5. 股票与指数对比

```bash
efi-cli stock compare 600519 沪深300 --beg 20240101 --end 20241231 --format table
efi-cli stock compare 600519 沪深300 --align-date --format json --pretty
```

### 6. 今日涨停板

```bash
efi-cli stock pool zt --format table
efi-cli stock pool zt --sort seal_amount --desc --limit 10
```

### 7. 板块概念热点

```bash
efi-cli stock sector list --type concept --sort pct --desc --limit 20 --format table
efi-cli stock sector members BK0477 --sort pct --desc --limit 10 --format table
```

### 8. 北向资金流入

```bash
efi-cli stock connect history --type north --limit 10 --format table
efi-cli stock connect hold-rank --type north --limit 20 --format table
```

### 9. 融资融券余额趋势

```bash
efi-cli stock margin account --limit 10 --format table
```

### 10. 限售解禁日历

```bash
efi-cli stock restricted summary --start 2026-04-01 --end 2026-06-30 --format table
efi-cli stock restricted detail --start 2026-04-01 --end 2026-04-30 --format table
```

### 11. 基金筛选

```bash
efi-cli fund screen --type gp --size-min 10 --1y-min 10 --3y-min 30 --sort 3y --limit 20
```

## 当前限制

1. 部分 `pool` 子命令（`strong`、`sub-new`、`zt-prev`）依赖的 `push2ex` 接口可能需要会话认证，偶尔返回空数据。
2. `bond` 和 `futures` 的 `history` 目前需要直接使用 `secid`。
3. 实时数据依赖上游接口可用性，高频请求可能触发限流（EOF 错误）。
4. `stock screen` 首次请求最多拉取 10000 只股票进行本地过滤。
5. `fund screen --sort 5y` 需要额外调用阶段收益 API 补充 5 年数据，查询较慢。
6. `stock connect hold` 数据源为按机构持股明细（非个股汇总），且历史数据仅到 2024-09，后续需更换 API。

## 推荐实践

1. 查询歧义标的前，先执行 `efi-cli search`。
2. 在 AI agent 中优先使用 `--format json --pretty` 或 `--compact`。
3. 新接入命令前先执行 `--schema` 和 `--list-fields` 了解能力。
4. 排查问题时优先看 `--raw` 输出。
5. 大批量查询建议适当控制并发，避免触发限流。
6. 使用 `--limit` 控制输出量，避免终端刷屏。
