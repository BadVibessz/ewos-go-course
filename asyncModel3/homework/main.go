package main

import (
	"context"
	"fmt"
	"hw-async/generator"
	"hw-async/pipeline"
	"hw-async/utils"
	"math"
	"slices"
	"time"

	log "github.com/sirupsen/logrus"
	dom "hw-async/domain"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}
var logger = log.New()

const (
	output1m  = "candles_1m_log.csv"
	output2m  = "candles_2m_log.csv"
	output10m = "candles_10m_log.csv"

	outputPerm = 0o664
)

func getClosedCandleIdx(candles []dom.Candle) int {
	for i, c := range candles {
		if c.Close != -1 {
			return i
		}
	}
	return -1
}

func CandleStage1m(ctx context.Context, pricesChan <-chan dom.Price, _ <-chan dom.Candle) <-chan dom.Candle {
	out := make(chan dom.Candle)

	count := 0
	resultLog := ""

	go func() {
		defer close(out)

		candles := make([]dom.Candle, 0, len(tickers))

		for p := range pricesChan {
			select {
			case <-ctx.Done():
			// todo:

			default:
				ts, err := dom.PeriodTS(dom.CandlePeriod1m, p.TS)
				if err == nil {
					idx := slices.IndexFunc(candles, func(candle dom.Candle) bool { return candle.Ticker == p.Ticker })
					if idx != -1 {
						candle := &candles[idx]

						candle.Low = math.Min(candle.Low, p.Value)
						candle.High = math.Max(candle.High, p.Value)

						if ts == candle.TS.Add(1*time.Minute) {
							candle.Close = p.Value
						}

					} else {
						candles = append(candles, dom.Candle{
							Ticker: p.Ticker,
							Period: dom.CandlePeriod1m,
							Open:   p.Value,
							High:   p.Value,
							Low:    p.Value,
							Close:  -1,
							TS:     ts,
						})
					}
				}

				idx := getClosedCandleIdx(candles)
				if idx != -1 {
					closed := candles[idx]

					logger.Infof("1m candle created: %+v\n", closed)
					resultLog += fmt.Sprintf("%+v\n", closed)

					out <- closed // send proceeded candle to channel

					candles = slices.Delete(candles, idx, idx+1)
					count++
				}

				if count == len(tickers) {
					err = utils.WriteToFile(output1m, []byte(resultLog), outputPerm) // todo: serialize to csv!!
					if err != nil {
						logger.Errorf("cannot save 1m candles to file: %v+\n", err)
					}

					return
				}
			}
		}
	}()

	return out
}

//func candleStage(ctx context.Context, period dom.CandlePeriod, pricesChan <-chan dom.Price, candleChan <-chan dom.Candle) <-chan dom.Candle {
//	out := make(chan dom.Candle)
//
//	go func() {
//		defer close(out)
//
//		n, err := utils.CandlePeriodToInt(period)
//		if err != nil {
//			return // todo
//		}
//
//		candles := make([]dom.Candle, 0, len(tickers))
//
//		count := 0
//
//		resultLog := ""
//
//		for price := range pricesChan {
//			select {
//			case <-ctx.Done():
//			// todo:
//
//			default:
//				c, open := <-candleChan
//				if open {
//					idx := slices.IndexFunc(candles, func(candle dom.Candle) bool { return c.Ticker == candle.Ticker })
//					if idx != -1 {
//						candle := &candles[idx]
//
//						ts, err := dom.PeriodTS(period, c.TS)
//						if err == nil {
//							if ts.Add(time.Duration(n)*time.Minute) == candle.TS {
//								candle.Close = price.Value
//							} else {
//								candle.Close = -1
//							}
//
//							candle.TS = ts
//							candle.Open = c.Open
//
//							candle.Low = math.Min(candle.Low, c.Low)
//							candle.High = math.Max(candle.High, c.High)
//						}
//
//					} else {
//						c.Close = -1
//						c.Period = period
//
//						candles = append(candles, c)
//					}
//				}
//
//				ts, err := dom.PeriodTS(period, price.TS)
//				if err == nil {
//					idx := slices.IndexFunc(candles, func(candle dom.Candle) bool { return candle.Ticker == price.Ticker })
//					if idx != -1 {
//						candle := &candles[idx]
//
//						candle.Low = math.Min(candle.Low, price.Value)
//						candle.High = math.Max(candle.High, price.Value)
//
//						//if ts == candle.TS.Add(time.Duration(n)*time.Minute) {
//						//if ts == candle.TS.Add(time.Duration(n)*time.Minute) || ts.Add(time.Duration(n)*time.Minute) == candle.TS {
//						if ts == candle.TS.Add(time.Duration(n)*time.Minute) && candle.Open != -1 {
//							if candle.Open == -1 {
//								logger.Errorf("aaaaaaaaaaaaaaaaaaaaaa") // TODO: UNDERSTAND WHY
//							}
//							candle.Close = price.Value
//						}
//
//					} else {
//						candles = append(candles, dom.Candle{
//							Ticker: price.Ticker,
//							Period: period,
//							High:   price.Value,
//							Low:    price.Value,
//							Open:   -1, //todo: init with price.Value for 1m candle // todo: open not copies from existing candle by some reason
//							Close:  -1,
//							TS:     ts,
//						})
//					}
//				}
//
//				idx := getClosedCandleIdx(candles)
//				if idx != -1 {
//					closed := candles[idx]
//
//					logger.Infof("%s candle created: %+v\n", period, closed)
//					resultLog += fmt.Sprintf("%+v\n", closed)
//
//					out <- closed // send proceeded candle to channel
//
//					candles = slices.Delete(candles, idx, idx+1) // todo: encapsulate (utils)
//					count++
//				}
//
//				if count == len(tickers) {
//					return
//				}
//
//			}
//		}
//	}()
//
//	return out
//}

func candleStage(ctx context.Context, period dom.CandlePeriod, pricesChan <-chan dom.Price, candleChan <-chan dom.Candle) <-chan dom.Candle {
	out := make(chan dom.Candle)

	go func() {
		defer close(out)

		n, err := utils.CandlePeriodToInt(period)
		if err != nil {
			return // todo
		}

		candles := make([]dom.Candle, 0, len(tickers))

		count := 0

		resultLog := ""

		select {
		case <-ctx.Done():
		// todo:

		case price, open := <-pricesChan:
			if open {
				ts, err := dom.PeriodTS(period, price.TS)
				if err == nil {
					idx := slices.IndexFunc(candles, func(candle dom.Candle) bool { return candle.Ticker == price.Ticker })
					if idx != -1 {
						candle := &candles[idx]

						candle.Low = math.Min(candle.Low, price.Value)
						candle.High = math.Max(candle.High, price.Value)

						//if ts == candle.TS.Add(time.Duration(n)*time.Minute) {
						//if ts == candle.TS.Add(time.Duration(n)*time.Minute) || ts.Add(time.Duration(n)*time.Minute) == candle.TS {
						if ts == candle.TS.Add(time.Duration(n)*time.Minute) && candle.Open != -1 {
							if candle.Open == -1 {
								logger.Errorf("aaaaaaaaaaaaaaaaaaaaaa") // TODO: UNDERSTAND WHY
							}
							candle.Close = price.Value
						}

					} else {
						candles = append(candles, dom.Candle{
							Ticker: price.Ticker,
							Period: period,
							High:   price.Value,
							Low:    price.Value,
							Open:   -1, //todo: init with price.Value for 1m candle // todo: open not copies from existing candle by some reason
							Close:  -1,
							TS:     ts,
						})
					}
				}
			}

		case c, open := <-candleChan:
			if open {
				idx := slices.IndexFunc(candles, func(candle dom.Candle) bool { return c.Ticker == candle.Ticker })
				if idx != -1 {
					candle := &candles[idx]

					ts, err := dom.PeriodTS(period, c.TS)
					if err == nil {
						if ts.Add(time.Duration(n)*time.Minute) == candle.TS {
							candle.Close = price.Value // todo: embedded Candle struct with bool filed 'Closed' ??
						} else {
							candle.Close = -1
						}

						candle.TS = ts
						candle.Open = c.Open

						candle.Low = math.Min(candle.Low, c.Low)
						candle.High = math.Max(candle.High, c.High)
					}

				} else {
					c.Close = -1
					c.Period = period

					candles = append(candles, c)
				}
			}

			idx := getClosedCandleIdx(candles)
			if idx != -1 {
				closed := candles[idx]

				logger.Infof("%s candle created: %+v\n", period, closed)
				resultLog += fmt.Sprintf("%+v\n", closed)

				out <- closed // send proceeded candle to channel

				candles = slices.Delete(candles, idx, idx+1) // todo: encapsulate (utils)
				count++
			}

			if count == len(tickers) {
				return
			}

		}
	}()

	return out
}

//func candleStagee(ctx context.Context, period dom.CandlePeriod, pricesChan <-chan dom.Price, candleChan <-chan dom.Candle) <-chan dom.Candle {
//	out := make(chan dom.Candle)
//
//	go func() {
//		defer close(out)
//
//		n, err := utils.CandlePeriodToInt(period)
//		if err != nil {
//			return // todo
//		}
//
//		candles := make([]dom.Candle, 0, len(tickers))
//
//		count := 0
//
//		resultLog := ""
//
//		go func() {
//			for price := range pricesChan {
//				ts, err := dom.PeriodTS(period, price.TS)
//				if err == nil {
//					idx := slices.IndexFunc(candles, func(candle dom.Candle) bool { return candle.Ticker == price.Ticker })
//					if idx != -1 {
//						candle := &candles[idx]
//
//						candle.Low = math.Min(candle.Low, price.Value)
//						candle.High = math.Max(candle.High, price.Value)
//
//						//if ts == candle.TS.Add(time.Duration(n)*time.Minute) {
//						if ts == candle.TS.Add(time.Duration(n)*time.Minute) || ts.Add(time.Duration(n)*time.Minute) == candle.TS {
//							if candle.Open == -1 {
//								logger.Errorf("aaaaaaaaaaaaaaaaaaaaaa") // TODO: UNDERSTAND WHY
//							}
//							candle.Close = price.Value
//						}
//
//					} else {
//						candles = append(candles, dom.Candle{
//							Ticker: price.Ticker,
//							Period: period,
//							High:   price.Value,
//							Low:    price.Value,
//							Open:   -1, //todo: init with price.Value for 1m candle // todo: open not copies from existing candle by some reason
//							Close:  -1,
//							TS:     ts,
//						})
//					}
//				}
//			}
//		}()
//
//		go func() {
//			for c := range candleChan {
//				idx := slices.IndexFunc(candles, func(candle dom.Candle) bool { return c.Ticker == candle.Ticker })
//				if idx != -1 {
//					candle := &candles[idx]
//
//					ts, err := dom.PeriodTS(period, c.TS)
//					if err == nil {
//						if ts == candle.TS.Add(time.Duration(n)*time.Minute) || ts.Add(time.Duration(n)*time.Minute) == candle.TS {
//							candle.Close = price.Value
//						} else {
//							candle.Close = -1
//						}
//
//						candle.TS = ts
//						candle.Open = c.Open
//
//						candle.Low = math.Min(candle.Low, c.Low)
//						candle.High = math.Max(candle.High, c.High)
//					}
//
//				} else {
//					c.Close = -1
//					c.Period = period
//
//					candles = append(candles, c)
//				}
//			}
//		}()
//
//		for {
//			select {
//			case <-ctx.Done():
//				// todo
//			}
//		}
//
//		for price := range pricesChan {
//			select {
//			case <-ctx.Done():
//			// todo:
//
//			default:
//				c, open := <-candleChan
//				if open {
//
//				}
//
//			}
//
//			idx := getClosedCandleIdx(candles)
//			if idx != -1 {
//				closed := candles[idx]
//
//				logger.Infof("%s candle created: %+v\n", period, closed)
//				resultLog += fmt.Sprintf("%+v\n", closed)
//
//				out <- closed // send proceeded candle to channel
//
//				candles = slices.Delete(candles, idx, idx+1) // todo: encapsulate (utils)
//				count++
//			}
//
//			if count == len(tickers) {
//				return
//			}
//
//		}
//	}
//}()

// return out
// }

// todo:!!!!
//func CandleStage1mm(ctx context.Context, pricesChan <-chan dom.Price, _ <-chan dom.Candle) <-chan dom.Candle {
//	return candleStage(ctx, pricesChan, _, )
//}

func CandleStage2m(ctx context.Context, pricesChan <-chan dom.Price, candleChan <-chan dom.Candle) <-chan dom.Candle {
	return candleStage(ctx, dom.CandlePeriod2m, pricesChan, candleChan)
}

func CandleStage3m(ctx context.Context, pricesChan <-chan dom.Price, candleChan <-chan dom.Candle) <-chan dom.Candle {
	return candleStage(ctx, "3m", pricesChan, candleChan)
}

func CandleStage10m(ctx context.Context, pricesChan <-chan dom.Price, candleChan <-chan dom.Candle) <-chan dom.Candle {
	return candleStage(ctx, dom.CandlePeriod10m, pricesChan, candleChan)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 5,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")
	prices := pg.Prices(ctx)

	candles := pipeline.ExecutePipeline(ctx, prices, CandleStage1m, CandleStage2m)

	//for c := range candles {
	//	logger.Infof("%+v\n", c)
	//}

	buf := make([]dom.Candle, 0)
	for c := range candles {
		buf = append(buf, c)
	}

	for _, c := range buf {
		logger.Infof("%+v\n", c)
	}

	cancel()
}
