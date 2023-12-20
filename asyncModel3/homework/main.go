package main

import (
	"context"
	"fmt"
	"hw-async/generator"
	"hw-async/pipeline"
	"hw-async/utils"
	"math"
	"os"
	"os/signal"
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

func pricesToCandles(pricesChan <-chan dom.Price) <-chan dom.Candle {
	out := make(chan dom.Candle)

	go func() {
		defer close(out)

		for p := range pricesChan {
			logger.Infof("PRICE: %+v", p)

			out <- utils.PriceToCandle(p)
		}
	}()

	return out
}

func candleStage(period dom.CandlePeriod, candleChan <-chan dom.Candle, errChan chan error) <-chan dom.Candle {
	out := make(chan dom.Candle)

	go func() {
		defer close(out)

		candles := make(map[string]dom.Candle)

		for c := range candleChan {
			candle, exist := candles[c.Ticker]

			ts, err := dom.PeriodTS(period, c.TS)
			if err != nil {
				errChan <- err
				continue
			}

			if exist { // not proceeded candle exists
				candle.Close = c.Close
				candle.Low = math.Min(candle.Low, c.Low)
				candle.High = math.Max(candle.High, c.High)

				if ts == candle.TS {
					out <- candle

					// delete proceeded ticker
					delete(candles, candle.Ticker)

					logger.Infof("%s candle created: %+v\n", period, candle)
				}

			} else {
				c.Close = -1
				c.Period = period
				c.TS = ts

				candles[c.Ticker] = c
			}
		}
	}()

	return out
}

func CandleStage1m(candleChan <-chan dom.Candle, errChan chan error) <-chan dom.Candle {
	return candleStage(dom.CandlePeriod1m, candleChan, errChan)
}

func CandleStage2m(candleChan <-chan dom.Candle, errChan chan error) <-chan dom.Candle {
	return candleStage(dom.CandlePeriod2m, candleChan, errChan)
}

func CandleStage10m(candleChan <-chan dom.Candle, errChan chan error) <-chan dom.Candle {
	return candleStage(dom.CandlePeriod10m, candleChan, errChan)
}

func writeCandlesToCsv(candles []dom.Candle, path string) error {
	output := ""
	for _, candle := range candles {
		output += fmt.Sprintf("%s\n", utils.CandleToCsv(candle))
	}

	err := utils.WriteToFile(path, []byte(output), outputPerm)
	if err != nil {
		return err
	}

	return nil
}

func LogCandlesToFile(path string, batchSize int) pipeline.Stage {
	return func(inCandles pipeline.InCandles, errChan pipeline.InErrors) pipeline.OutCandles {
		out := make(chan dom.Candle)

		candles := make([]dom.Candle, 0, batchSize)

		go func() {
			defer close(out)

			for c := range inCandles {
				candles = append(candles, c)

				if len(candles) == batchSize {

					err := writeCandlesToCsv(candles, path)
					if err != nil {
						errChan <- err
					}

					// clear slice
					candles = make([]dom.Candle, 0, batchSize)
				}

				// return read candle
				out <- c
			}

			// if chan is closed and there's candles in slice => output
			if len(candles) != 0 {
				err := writeCandlesToCsv(candles, path)
				if err != nil {
					errChan <- err
				}
			}
		}()

		return out
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

	// register signal interrupt handler
	go func() {
		for {
			sig := <-sigChan

			if sig == os.Interrupt {
				logger.Info("Interrupting program...")
				cancel() // cancel context
			}
		}
	}()

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")
	prices := pg.Prices(ctx)

	candleChan := pricesToCandles(prices) // TODO: STOPS WORKING
	candles, errs := pipeline.ExecutePipeline(ctx, candleChan, CandleStage1m, LogCandlesToFile(output1m, 100),
		CandleStage2m, LogCandlesToFile(output2m, 1),
		CandleStage10m, LogCandlesToFile(output10m, 1))

	for err := range errs {
		logger.Error(err.Error())
	}

	buf := make([]dom.Candle, 0)
	for c := range candles {
		buf = append(buf, c)
	}

	for _, c := range buf {
		logger.Infof("%+v\n", c)
	}

	cancel()
}
