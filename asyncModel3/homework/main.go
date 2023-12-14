package main

import (
	"context"
	"fmt"
	"hw-async/generator"
	"hw-async/pipeline"
	"math"
	"slices"
	"time"

	log "github.com/sirupsen/logrus"
	dom "hw-async/domain"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

func areCandlesClosed(candles []dom.Candle) bool {
	for _, c := range candles {
		if c.Close == -1 {
			return false
		}
	}
	return true
}

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

	fmt.Println("STARTING GENERATING 1M CANDLES")

	count := 0

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
					out <- candles[idx] // send proceeded candle to channel

					candles = slices.Delete(candles, idx, idx+1) // todo: encapsulate (utils)
					count++
				}

				if count == len(tickers) {
					return
				}
			}
		}
	}()

	return out
}

func CandleStage2m(ctx context.Context, pricesChan <-chan dom.Price, candleChan <-chan dom.Candle) <-chan dom.Candle {
	out := make(chan dom.Candle)

	fmt.Println("STARTING GENERATING 2M CANDLES")

	go func() {
		defer close(out)

		candles := make([]dom.Candle, 0, len(tickers))

		count := 0

		for price := range pricesChan {
			select {
			case <-ctx.Done():
			// todo:

			default:
				c, open := <-candleChan
				if !open {
					return // todo
				}

				idx := slices.IndexFunc(candles, func(c dom.Candle) bool { return c.Ticker == c.Ticker })
				if idx != -1 {
					candle := &candles[idx]

					candle.TS = c.TS
					candle.Open = c.Open
					candle.Close = -1

					candle.Low = math.Min(candle.Low, c.Low)
					candle.High = math.Max(candle.High, c.High)
				} else {
					c.Close = -1
					c.Period = dom.CandlePeriod2m

					candles = append(candles, c)
				}

				ts, err := dom.PeriodTS(dom.CandlePeriod2m, price.TS)
				if err == nil {
					idx = slices.IndexFunc(candles, func(candle dom.Candle) bool { return candle.Ticker == price.Ticker })
					if idx != -1 {
						candle := &candles[idx]

						candle.Low = math.Min(candle.Low, price.Value)
						candle.High = math.Max(candle.High, price.Value)

						if ts == candle.TS.Add(2*time.Minute) {
							candle.Close = price.Value
						}

					} else {
						candles = append(candles, dom.Candle{
							Ticker: price.Ticker,
							Period: dom.CandlePeriod2m,
							High:   price.Value,
							Low:    price.Value,
						})
					}
				}

				idx = getClosedCandleIdx(candles)
				if idx != -1 {
					out <- candles[idx] // send proceeded candle to channel

					candles = slices.Delete(candles, idx, idx+1) // todo: encapsulate (utils)
					count++
				}

				if count == len(tickers) {
					return
				}

			}
		}
	}()

	return out
}

func main() {
	logger := log.New()
	ctx, cancel := context.WithCancel(context.Background())

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")
	prices := pg.Prices(ctx)

	candles := pipeline.ExecutePipeline(ctx, prices, CandleStage1m, CandleStage2m)

	for c := range candles {
		logger.Infof("%+v\n", c)
	}

	cancel()
}
