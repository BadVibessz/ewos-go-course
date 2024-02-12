package main

import (
	"context"
	"fmt"
	"hw-async/generator"
	"hw-async/graphic"
	"hw-async/pipeline"
	"hw-async/utils"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	dom "hw-async/domain"
)

var (
	tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

	logger = log.New()
)

const (
	log1mPath  = "output/logs/candles_1m_log.csv"
	log2mPath  = "output/logs/candles_2m_log.csv"
	log10mPath = "output/logs/candles_10m_log.csv"

	graph1mPath  = "output/graphs/1m.txt"
	graph2mPath  = "output/graphs/2m.txt"
	graph10mPath = "output/graphs/10m.txt"

	writePerm = 0o664
)

const (
	batchSize = 10
	genFactor = 10
)

func pricesToCandles(pricesChan <-chan dom.Price) <-chan dom.Candle {
	out := make(chan dom.Candle)

	go func() {
		defer close(out)

		for p := range pricesChan {
			out <- utils.PriceToCandle(p)
		}
	}()

	return out
}

func candleStage(period dom.CandlePeriod, candleChan <-chan dom.Candle, errChan chan error) <-chan dom.Candle {
	out := make(chan dom.Candle)

	go func() {
		defer close(out)

		n, err := utils.CandlePeriodToInt(period)
		if err != nil {
			errChan <- err
			return
		}

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
				candle.Low = utils.Min(candle.Open, candle.Low, candle.High, candle.Close, c.Open, c.Low, c.High, c.Close)
				candle.High = utils.Max(candle.Open, candle.Low, candle.High, candle.Close, c.Open, c.Low, c.High, c.Close)

				if ts == candle.TS.Add(time.Duration(n)*time.Minute) {
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

	err := utils.WriteToFile(path, []byte(output), writePerm)
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

func LogCandlesGraphToFile(path string) pipeline.Stage {
	return func(inCandles pipeline.InCandles, errChan pipeline.InErrors) pipeline.OutCandles {
		out := make(chan dom.Candle)

		candles := make(map[string][]dom.Candle)

		go func() {
			defer close(out)

			for c := range inCandles {
				candles[c.Ticker] = append(candles[c.Ticker], c)

				// return read candle
				out <- c
			}

			// if chan is closed and there's candles in slice => output graph
			if len(candles) != 0 {
				for t, s := range candles {
					timingPrices := make([]graphic.TimePrice, 0, len(candles))

					for _, c := range s {
						timingPrices = append(timingPrices, graphic.TimePrice{First: c.TS, Second: c.Close})
					}

					graph := graphic.New(t, timingPrices...)

					resLog := graph.GenerateString('*')

					err := utils.WriteToFile(path, []byte(resLog), writePerm)
					if err != nil {
						errChan <- err
						return
					}
				}
			}
		}()

		return out
	}
}

func clearAllFiles(paths ...string) error {
	for _, p := range paths {
		err := utils.ClearFile(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	err := clearAllFiles(log1mPath, log2mPath, log10mPath, graph1mPath, graph2mPath, graph10mPath)
	if err != nil {
		logger.Fatalf("error occuerd while clearing files: %s", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

	// register signal interrupt handler
	go func() {
		for {
			sig := <-sigChan

			if sig == os.Interrupt {
				logger.Info("interrupting program...")
				cancel() // cancel context
			}
		}
	}()

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  genFactor,
		Delay:   time.Millisecond * 1,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")

	prices := pg.Prices(ctx)
	candleChan := pricesToCandles(prices)

	candles, errs := pipeline.ExecutePipeline(candleChan,
		CandleStage1m, LogCandlesToFile(log1mPath, batchSize), LogCandlesGraphToFile(graph1mPath),
		CandleStage2m, LogCandlesToFile(log2mPath, batchSize), LogCandlesGraphToFile(graph2mPath),
		CandleStage10m, LogCandlesToFile(log10mPath, batchSize), LogCandlesGraphToFile(graph10mPath),
	)

	// log and handle errors occurred in pipeline stages

	// listen for occurred errors
	go func() {
		for err = range errs {
			logger.Error(err.Error())
		}
	}()

	// read from chan to make it work
	for c := range candles {
		logger.Infof("pipeline output: %+v\n", c)
	}
}
