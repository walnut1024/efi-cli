package stock

import (
	"github.com/spf13/cobra"
	"github.com/walnut1024/efi-cli/internal/analytics"
	"github.com/walnut1024/efi-cli/internal/cliutil"
	"github.com/walnut1024/efi-cli/internal/model"
)

var (
	stockKlt              int
	stockFqt              int
	stockBeg              string
	stockEnd              string
	stockDate             string
	stockMaxCount         int
	stockScreenPEMin      float64
	stockScreenPEMax      float64
	stockScreenPBMin      float64
	stockScreenPBMax      float64
	stockScreenROEMin     float64
	stockScreenCapMin     float64
	stockScreenCapMax     float64
	stockScreenPctMin     float64
	stockScreenPctMax     float64
	stockScreenMarket     string
	stockFinanceType      string
	stockIndicators       []string
	stockStats            []string
	stockSummary          bool
	stockCompareAlign     bool
	stockCompareMetric    string
	stockSectorType       string
	stockConnectType      string
	stockConnectIndicator string
)

func formatOrExit(records []model.Record) { stockOutput.FormatOrExit(records) }

type recordFetcher func(args []string) ([]model.Record, error)
type rawFetcher func(args []string) ([]byte, error)
type validator func() error

type commandRunOptions struct {
	schema   *cliutil.CommandSchema
	validate validator
	raw      rawFetcher
	records  recordFetcher
	sort     bool
}

func runFormattedCommand(args []string, fetch recordFetcher) {
	records, err := fetch(args)
	if err != nil {
		cliutil.PrintErrorAndExit(err, stockOutput)
	}
	formatOrExit(records)
}

func runCommand(args []string, opts commandRunOptions) {
	if opts.validate != nil {
		if err := opts.validate(); err != nil {
			cliutil.PrintErrorAndExit(err, stockOutput)
		}
	}
	if stockOutput.HandleSchemaOrListFields(opts.schema) {
		return
	}
	if stockOutput.Raw && opts.raw != nil {
		raw, err := opts.raw(args)
		if err != nil {
			cliutil.PrintErrorAndExit(err, stockOutput)
		}
		stockOutput.PrintRawOrExit(raw)
		return
	}
	records, err := opts.records(args)
	if err != nil {
		cliutil.PrintErrorAndExit(err, stockOutput)
	}
	if opts.sort {
		records, err = stockOutput.ApplySort(records)
		if err != nil {
			cliutil.PrintErrorAndExit(err, stockOutput)
		}
	}
	stockOutput.FormatWithSchemaOrExit(records, opts.schema)
}

func addCommandOutputFlags(cmd *cobra.Command, inspect bool, sortable bool) {
	stockOutput.AddFlags(cmd)
	if inspect {
		stockOutput.AddInspectFlags(cmd)
	}
	if sortable {
		stockOutput.AddSortFlags(cmd)
	}
}

func AddCommands(root *cobra.Command) {
	stockCmd := &cobra.Command{
		Use:   "stock",
		Short: "Stock market data commands",
	}
	root.AddCommand(stockCmd)

	// efi quote <codes...>
	quoteCmd := &cobra.Command{
		Use:   "quote [codes]",
		Short: "Get realtime quote for stocks",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema:   &quoteSchema,
				validate: stockOutput.ValidateSort,
				raw:      GetQuoteRaw,
				records:  GetQuote,
				sort:     true,
			})
		},
	}
	addCommandOutputFlags(quoteCmd, true, true)
	stockCmd.AddCommand(quoteCmd)

	// efi history <code>
	historyCmd := &cobra.Command{
		Use:   "history [code]",
		Short: "Get K-line history (default: daily, pre-adjusted)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema: &historySchema,
				raw: func(args []string) ([]byte, error) {
					return GetHistoryRaw(args[0], stockBeg, stockEnd, stockKlt, stockFqt)
				},
				records: func(args []string) ([]model.Record, error) {
					return GetHistory(args[0], stockBeg, stockEnd, stockKlt, stockFqt, analytics.HistoryAnalysisOptions{
						Indicators: parseCSVFlagValues(stockIndicators),
						Stats:      parseCSVFlagValues(stockStats),
						Summary:    stockSummary,
					})
				},
			})
		},
	}
	historyCmd.Flags().IntVar(&stockKlt, "klt", 101, "Period: 1/5/15/30/60min, 101=day, 102=week, 103=month")
	historyCmd.Flags().IntVar(&stockFqt, "fqt", 1, "Adjust: 0=none, 1=forward, 2=backward")
	historyCmd.Flags().StringVar(&stockBeg, "beg", "19000101", "Start date YYYYMMDD")
	historyCmd.Flags().StringVar(&stockEnd, "end", "20500101", "End date YYYYMMDD")
	historyCmd.Flags().StringSliceVar(&stockIndicators, "indicators", nil, "Indicators to add: ma[:5,10,20], ema[:12,26], macd, rsi[:14], boll[:20,2]")
	historyCmd.Flags().StringSliceVar(&stockStats, "stats", nil, "Summary stats: total_return,cumulative_pct,amplitude_avg,high,low,max_drawdown,start_close,end_close,start_date,end_date,bars")
	historyCmd.Flags().BoolVar(&stockSummary, "summary", false, "Return a single summary row instead of full history rows")
	addCommandOutputFlags(historyCmd, true, false)
	stockCmd.AddCommand(historyCmd)

	compareCmd := &cobra.Command{
		Use:   "compare [left] [right]",
		Short: "Compare two stocks or indices over a date range",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema: &compareSchema,
				validate: func() error {
					if err := validateCompareMetric(stockCompareMetric); err != nil {
						return err
					}
					if stockCompareAlign {
						return nil
					}
					return stockOutput.ValidateSort()
				},
				records: func(args []string) ([]model.Record, error) {
					return CompareTargets(args[0], args[1], stockBeg, stockEnd, stockKlt, stockCompareAlign, stockCompareMetric)
				},
				sort: !stockCompareAlign,
			})
		},
	}
	compareCmd.Flags().StringVar(&stockBeg, "beg", "19000101", "Start date YYYYMMDD")
	compareCmd.Flags().StringVar(&stockEnd, "end", "20500101", "End date YYYYMMDD")
	compareCmd.Flags().IntVar(&stockKlt, "klt", 101, "Period: 101=day, 102=week, 103=month")
	compareCmd.Flags().BoolVar(&stockCompareAlign, "align-date", false, "Align output by common dates and emit side-by-side series")
	compareCmd.Flags().StringVar(&stockCompareMetric, "metric", "", "Metric used for summary comparison: total_return, max_drawdown, high, low, amplitude_avg, start_close, end_close, bars")
	addCommandOutputFlags(compareCmd, true, true)
	stockCmd.AddCommand(compareCmd)

	// efi realtime [market]
	realtimeCmd := &cobra.Command{
		Use:   "realtime [market]",
		Short: "Get market-wide realtime quotes (default: 沪深京A股)",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema:   &realtimeSchema,
				validate: stockOutput.ValidateSort,
				raw: func(args []string) ([]byte, error) {
					market := ""
					if len(args) > 0 {
						market = args[0]
					}
					return GetRealtimeRaw(market)
				},
				records: func(args []string) ([]model.Record, error) {
					market := ""
					if len(args) > 0 {
						market = args[0]
					}
					return GetRealtime(market)
				},
				sort: true,
			})
		},
	}
	addCommandOutputFlags(realtimeCmd, true, true)
	stockCmd.AddCommand(realtimeCmd)

	// efi bill <code>
	billCmd := &cobra.Command{
		Use:   "bill [code]",
		Short: "Get historical capital flow",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetHistoryBill(args[0])
			})
		},
	}
	addCommandOutputFlags(billCmd, false, false)
	stockCmd.AddCommand(billCmd)

	// efi todaybill <code>
	todayBillCmd := &cobra.Command{
		Use:   "todaybill [code]",
		Short: "Get intraday minute-level capital flow",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetTodayBill(args[0])
			})
		},
	}
	addCommandOutputFlags(todayBillCmd, false, false)
	stockCmd.AddCommand(todayBillCmd)

	// efi deal <code>
	dealCmd := &cobra.Command{
		Use:   "deal [code]",
		Short: "Get latest trading day deal details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetDealDetail(args[0], stockMaxCount)
			})
		},
	}
	dealCmd.Flags().IntVar(&stockMaxCount, "max", 10000, "Max records")
	addCommandOutputFlags(dealCmd, false, false)
	stockCmd.AddCommand(dealCmd)

	// efi info <code>
	infoCmd := &cobra.Command{
		Use:   "info [code]",
		Short: "Get stock base info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				record, err := GetBaseInfo(args[0])
				if err != nil {
					return nil, err
				}
				return []model.Record{record}, nil
			})
		},
	}
	addCommandOutputFlags(infoCmd, false, false)
	stockCmd.AddCommand(infoCmd)

	// efi billboard --start <date> --end <date>
	billboardCmd := &cobra.Command{
		Use:   "billboard",
		Short: "Get daily billboard (龙虎榜) data",
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetBillboard(stockBeg, stockEnd)
			})
		},
	}
	billboardCmd.Flags().StringVar(&stockBeg, "start", "", "Start date YYYY-MM-DD")
	billboardCmd.Flags().StringVar(&stockEnd, "end", "", "End date YYYY-MM-DD")
	addCommandOutputFlags(billboardCmd, false, false)
	stockCmd.AddCommand(billboardCmd)

	// efi board <code>
	boardCmd := &cobra.Command{
		Use:   "board [code]",
		Short: "Get boards a stock belongs to",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetBelongBoard(args[0])
			})
		},
	}
	addCommandOutputFlags(boardCmd, false, false)
	stockCmd.AddCommand(boardCmd)

	// efi members <index>
	membersCmd := &cobra.Command{
		Use:   "members [index]",
		Short: "Get index constituent stocks",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetMembers(args[0])
			})
		},
	}
	addCommandOutputFlags(membersCmd, false, false)
	stockCmd.AddCommand(membersCmd)

	// efi ipo
	ipoCmd := &cobra.Command{
		Use:   "ipo",
		Short: "Get latest IPO review status",
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetIPOInfo()
			})
		},
	}
	addCommandOutputFlags(ipoCmd, false, false)
	stockCmd.AddCommand(ipoCmd)

	// efi holders [--date <date>]
	holdersCmd := &cobra.Command{
		Use:   "holders",
		Short: "Get latest shareholder number changes",
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetHolderNumber(stockDate)
			})
		},
	}
	holdersCmd.Flags().StringVar(&stockDate, "date", "", "Report date YYYY-MM-DD")
	addCommandOutputFlags(holdersCmd, false, false)
	stockCmd.AddCommand(holdersCmd)

	// efi performance [--date <date>]
	perfCmd := &cobra.Command{
		Use:   "performance",
		Short: "Get company quarterly performance",
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetPerformance(stockDate)
			})
		},
	}
	perfCmd.Flags().StringVar(&stockDate, "date", "", "Report date YYYY-MM-DD")
	addCommandOutputFlags(perfCmd, false, false)
	stockCmd.AddCommand(perfCmd)

	// efi stock screen
	screenCmd := &cobra.Command{
		Use:   "screen",
		Short: "Screen stocks by PE/PB/ROE/market cap/pct",
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetScreen(stockScreenMarket, ScreenFilter{
					PEMin:  stockScreenPEMin,
					PEMax:  stockScreenPEMax,
					PBMin:  stockScreenPBMin,
					PBMax:  stockScreenPBMax,
					ROEMin: stockScreenROEMin,
					CapMin: stockScreenCapMin,
					CapMax: stockScreenCapMax,
					PctMin: stockScreenPctMin,
					PctMax: stockScreenPctMax,
				}, stockOutput.Limit)
			})
		},
	}
	screenCmd.Flags().Float64Var(&stockScreenPEMin, "pe-min", 0, "Min PE ratio")
	screenCmd.Flags().Float64Var(&stockScreenPEMax, "pe-max", 0, "Max PE ratio (0=unlimited)")
	screenCmd.Flags().Float64Var(&stockScreenPBMin, "pb-min", 0, "Min PB ratio")
	screenCmd.Flags().Float64Var(&stockScreenPBMax, "pb-max", 0, "Max PB ratio (0=unlimited)")
	screenCmd.Flags().Float64Var(&stockScreenROEMin, "roe-min", 0, "Min ROE (%)")
	screenCmd.Flags().Float64Var(&stockScreenCapMin, "cap-min", 0, "Min market cap (亿元)")
	screenCmd.Flags().Float64Var(&stockScreenCapMax, "cap-max", 0, "Max market cap (亿元, 0=unlimited)")
	screenCmd.Flags().Float64Var(&stockScreenPctMin, "pct-min", 0, "Min daily change (%)")
	screenCmd.Flags().Float64Var(&stockScreenPctMax, "pct-max", 0, "Max daily change (%)")
	screenCmd.Flags().StringVar(&stockScreenMarket, "market", "沪深京A股", "Market filter (沪深A股, 创业板, 科创板, etc.)")
	addCommandOutputFlags(screenCmd, false, false)
	stockCmd.AddCommand(screenCmd)

	// efi stock finance <code>
	financeCmd := &cobra.Command{
		Use:   "finance [code]",
		Short: "Get financial statements (income/balance/cashflow)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetFinance(args[0], stockFinanceType, stockOutput.Limit)
			})
		},
	}
	financeCmd.Flags().StringVar(&stockFinanceType, "type", "income", "Statement type: income, balance, cashflow")
	addCommandOutputFlags(financeCmd, false, false)
	stockCmd.AddCommand(financeCmd)

	// efi stock dividend <code>
	dividendCmd := &cobra.Command{
		Use:   "dividend [code]",
		Short: "Get dividend history",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetDividend(args[0], stockOutput.Limit)
			})
		},
	}
	addCommandOutputFlags(dividendCmd, false, false)
	stockCmd.AddCommand(dividendCmd)

	// efi stock pool
	poolCmd := &cobra.Command{
		Use:   "pool",
		Short: "Limit-up/down pools (涨停/跌停/炸板/强势/次新)",
	}
	stockCmd.AddCommand(poolCmd)

	poolTypes := []struct {
		use   string
		short string
		pt    string
	}{
		{"zt", "Limit-up pool (涨停池)", "zt"},
		{"zt-prev", "Previous day limit-up (昨日涨停)", "zt-prev"},
		{"dt", "Limit-down pool (跌停池)", "dt"},
		{"zb", "Broken limit pool (炸板池)", "zb"},
		{"strong", "Strong stocks pool (强势股池)", "strong"},
		{"sub-new", "Sub-new stocks pool (次新股池)", "sub-new"},
	}
	for _, pt := range poolTypes {
		pt := pt
		cmd := &cobra.Command{
			Use:   pt.use,
			Short: pt.short,
			Run: func(cmd *cobra.Command, args []string) {
				runCommand(args, commandRunOptions{
					schema:   poolSchemaForType(pt.pt),
					validate: stockOutput.ValidateSort,
					raw: func(args []string) ([]byte, error) {
						return GetPoolRaw(pt.pt, stockDate)
					},
					records: func(args []string) ([]model.Record, error) {
						return GetPool(pt.pt, stockDate)
					},
					sort: true,
				})
			},
		}
		cmd.Flags().StringVar(&stockDate, "date", "", "Date YYYYMMDD (default: today)")
		addCommandOutputFlags(cmd, true, true)
		poolCmd.AddCommand(cmd)
	}

	// efi stock sector
	sectorCmd := &cobra.Command{
		Use:   "sector",
		Short: "Sector/industry/concept board data (板块深数据)",
	}
	stockCmd.AddCommand(sectorCmd)

	sectorListCmd := &cobra.Command{
		Use:   "list",
		Short: "List sectors with realtime quotes",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema:   &sectorListSchema,
				validate: stockOutput.ValidateSort,
				records: func(args []string) ([]model.Record, error) {
					return GetSectorList(stockSectorType)
				},
				sort: true,
			})
		},
	}
	sectorListCmd.Flags().StringVar(&stockSectorType, "type", "concept", "Sector type: industry, concept")
	addCommandOutputFlags(sectorListCmd, true, true)
	sectorCmd.AddCommand(sectorListCmd)

	sectorMembersCmd := &cobra.Command{
		Use:   "members [code]",
		Short: "Get sector member stocks",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema:   &sectorMembersSchema,
				validate: stockOutput.ValidateSort,
				records: func(args []string) ([]model.Record, error) {
					return GetSectorMembers(args[0], stockOutput.Limit)
				},
				sort: true,
			})
		},
	}
	addCommandOutputFlags(sectorMembersCmd, true, true)
	sectorCmd.AddCommand(sectorMembersCmd)

	sectorHistoryCmd := &cobra.Command{
		Use:   "history [code]",
		Short: "Get sector K-line history",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetSectorHistory(args[0], stockBeg, stockEnd, stockKlt, stockFqt)
			})
		},
	}
	sectorHistoryCmd.Flags().IntVar(&stockKlt, "klt", 101, "Period: 101=day, 102=week, 103=month")
	sectorHistoryCmd.Flags().IntVar(&stockFqt, "fqt", 1, "Adjust: 0=none, 1=forward")
	sectorHistoryCmd.Flags().StringVar(&stockBeg, "beg", "19000101", "Start date YYYYMMDD")
	sectorHistoryCmd.Flags().StringVar(&stockEnd, "end", "20500101", "End date YYYYMMDD")
	addCommandOutputFlags(sectorHistoryCmd, false, false)
	sectorCmd.AddCommand(sectorHistoryCmd)

	sectorQuoteCmd := &cobra.Command{
		Use:   "quote [code]",
		Short: "Get sector realtime quote",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetSectorQuote(args[0])
			})
		},
	}
	addCommandOutputFlags(sectorQuoteCmd, false, false)
	sectorCmd.AddCommand(sectorQuoteCmd)

	// efi stock connect
	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "Stock Connect data (沪深港通)",
	}
	stockCmd.AddCommand(connectCmd)

	connectSummaryCmd := &cobra.Command{
		Use:   "summary",
		Short: "Stock Connect quota summary",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema:  &connectSummarySchema,
				records: func(args []string) ([]model.Record, error) { return GetConnectSummary() },
			})
		},
	}
	addCommandOutputFlags(connectSummaryCmd, true, false)
	connectCmd.AddCommand(connectSummaryCmd)

	connectHistoryCmd := &cobra.Command{
		Use:   "history",
		Short: "Stock Connect historical fund flow",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema: &connectHistorySchema,
				records: func(args []string) ([]model.Record, error) {
					return GetConnectHistory(stockConnectType, stockOutput.Limit)
				},
			})
		},
	}
	connectHistoryCmd.Flags().StringVar(&stockConnectType, "type", "north", "Type: north, sh, sz, south, hshk, szhk")
	addCommandOutputFlags(connectHistoryCmd, true, false)
	connectCmd.AddCommand(connectHistoryCmd)

	connectRealtimeCmd := &cobra.Command{
		Use:   "realtime",
		Short: "Stock Connect minute-level fund flow",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema:  &connectRealtimeSchema,
				records: func(args []string) ([]model.Record, error) { return GetConnectRealtime() },
			})
		},
	}
	addCommandOutputFlags(connectRealtimeCmd, true, false)
	connectCmd.AddCommand(connectRealtimeCmd)

	connectHoldRankCmd := &cobra.Command{
		Use:   "hold-rank",
		Short: "Stock Connect holdings ranking",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema:   &connectHoldRankSchema,
				validate: stockOutput.ValidateSort,
				records: func(args []string) ([]model.Record, error) {
					return GetConnectHoldRank(stockConnectType, stockConnectIndicator, stockDate)
				},
				sort: true,
			})
		},
	}
	connectHoldRankCmd.Flags().StringVar(&stockConnectType, "type", "north", "Type: north, sh, sz")
	connectHoldRankCmd.Flags().StringVar(&stockConnectIndicator, "indicator", "", "Period: 1,3,5,10,M,Q,Y")
	connectHoldRankCmd.Flags().StringVar(&stockDate, "date", "", "Date YYYY-MM-DD")
	addCommandOutputFlags(connectHoldRankCmd, true, true)
	connectCmd.AddCommand(connectHoldRankCmd)

	connectHoldCmd := &cobra.Command{
		Use:   "hold [code]",
		Short: "Individual stock holding detail via Stock Connect",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetConnectHold(args[0], stockBeg, stockEnd)
			})
		},
	}
	connectHoldCmd.Flags().StringVar(&stockBeg, "start", "", "Start date YYYY-MM-DD")
	connectHoldCmd.Flags().StringVar(&stockEnd, "end", "", "End date YYYY-MM-DD")
	addCommandOutputFlags(connectHoldCmd, false, false)
	connectCmd.AddCommand(connectHoldCmd)

	connectAHCmd := &cobra.Command{
		Use:   "ah",
		Short: "A+H dual-listed stocks comparison",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema:   &connectAHSchema,
				validate: stockOutput.ValidateSort,
				records: func(args []string) ([]model.Record, error) {
					return GetConnectAH()
				},
				sort: true,
			})
		},
	}
	addCommandOutputFlags(connectAHCmd, true, true)
	connectCmd.AddCommand(connectAHCmd)

	// efi stock hot-rank
	hotRankCmd := &cobra.Command{
		Use:   "hot-rank",
		Short: "Stock popularity ranking (人气榜)",
	}
	stockCmd.AddCommand(hotRankCmd)

	hotRankListCmd := &cobra.Command{
		Use:   "list",
		Short: "Current popularity ranking",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema:   &hotRankSchema,
				validate: stockOutput.ValidateSort,
				records: func(args []string) ([]model.Record, error) {
					return GetHotRank(stockOutput.Limit)
				},
				sort: true,
			})
		},
	}
	addCommandOutputFlags(hotRankListCmd, true, true)
	hotRankCmd.AddCommand(hotRankListCmd)

	hotRankHistoryCmd := &cobra.Command{
		Use:   "history [code]",
		Short: "Historical popularity trend for a stock",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetHotRankHistory(args[0])
			})
		},
	}
	addCommandOutputFlags(hotRankHistoryCmd, false, false)
	hotRankCmd.AddCommand(hotRankHistoryCmd)

	hotRankKeywordCmd := &cobra.Command{
		Use:   "keyword [code]",
		Short: "Hot keywords for a stock",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetHotRankKeyword(args[0])
			})
		},
	}
	addCommandOutputFlags(hotRankKeywordCmd, false, false)
	hotRankCmd.AddCommand(hotRankKeywordCmd)

	// efi stock comment
	commentCmd := &cobra.Command{
		Use:   "comment",
		Short: "Stock commentary scores (千股千评)",
	}
	stockCmd.AddCommand(commentCmd)

	commentListCmd := &cobra.Command{
		Use:   "list",
		Short: "Market-wide commentary scores",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema:   &commentSchema,
				validate: stockOutput.ValidateSort,
				records: func(args []string) ([]model.Record, error) {
					return GetComment(stockOutput.Limit)
				},
				sort: true,
			})
		},
	}
	addCommandOutputFlags(commentListCmd, true, true)
	commentCmd.AddCommand(commentListCmd)

	commentInstitutionCmd := &cobra.Command{
		Use:   "institution [code]",
		Short: "Institutional participation for a stock",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetCommentInstitution(args[0])
			})
		},
	}
	addCommandOutputFlags(commentInstitutionCmd, false, false)
	commentCmd.AddCommand(commentInstitutionCmd)

	// efi stock margin
	marginCmd := &cobra.Command{
		Use:   "margin",
		Short: "Margin trading data (融资融券)",
	}
	stockCmd.AddCommand(marginCmd)

	marginAccountCmd := &cobra.Command{
		Use:   "account",
		Short: "Market-wide margin account statistics",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema: &marginSchema,
				records: func(args []string) ([]model.Record, error) {
					return GetMarginAccount(stockOutput.Limit)
				},
			})
		},
	}
	addCommandOutputFlags(marginAccountCmd, true, false)
	marginCmd.AddCommand(marginAccountCmd)

	// efi stock restricted
	restrictedCmd := &cobra.Command{
		Use:   "restricted",
		Short: "Restricted share release data (限售解禁)",
	}
	stockCmd.AddCommand(restrictedCmd)

	restrictedSummaryCmd := &cobra.Command{
		Use:   "summary",
		Short: "Market-wide release summary",
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetRestrictedSummary(stockBeg, stockEnd, "")
			})
		},
	}
	restrictedSummaryCmd.Flags().StringVar(&stockBeg, "start", "", "Start date YYYY-MM-DD")
	restrictedSummaryCmd.Flags().StringVar(&stockEnd, "end", "", "End date YYYY-MM-DD")
	addCommandOutputFlags(restrictedSummaryCmd, false, false)
	restrictedCmd.AddCommand(restrictedSummaryCmd)

	restrictedDetailCmd := &cobra.Command{
		Use:   "detail",
		Short: "Individual stock release detail",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(args, commandRunOptions{
				schema: &restrictedDetailSchema,
				records: func(args []string) ([]model.Record, error) {
					return GetRestrictedDetail(stockBeg, stockEnd)
				},
			})
		},
	}
	restrictedDetailCmd.Flags().StringVar(&stockBeg, "start", "", "Start date YYYY-MM-DD")
	restrictedDetailCmd.Flags().StringVar(&stockEnd, "end", "", "End date YYYY-MM-DD")
	addCommandOutputFlags(restrictedDetailCmd, true, false)
	restrictedCmd.AddCommand(restrictedDetailCmd)

	restrictedQueueCmd := &cobra.Command{
		Use:   "queue [code]",
		Short: "Upcoming release batches for a stock",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetRestrictedQueue(args[0])
			})
		},
	}
	addCommandOutputFlags(restrictedQueueCmd, false, false)
	restrictedCmd.AddCommand(restrictedQueueCmd)

	restrictedHoldersCmd := &cobra.Command{
		Use:   "holders [code]",
		Short: "Release holder details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runFormattedCommand(args, func(args []string) ([]model.Record, error) {
				return GetRestrictedHolders(args[0], stockDate)
			})
		},
	}
	restrictedHoldersCmd.Flags().StringVar(&stockDate, "date", "", "Date YYYY-MM-DD")
	addCommandOutputFlags(restrictedHoldersCmd, false, false)
	restrictedCmd.AddCommand(restrictedHoldersCmd)
}
