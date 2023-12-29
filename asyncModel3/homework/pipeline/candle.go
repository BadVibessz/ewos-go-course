package pipeline

import (
	dom "hw-async/domain"
)

type (
	InCandles = <-chan dom.Candle
	InErrors  = chan error

	OutCandles = InCandles
	OutErrors  = InErrors
)

type Stage func(InCandles, InErrors) OutCandles

func ExecutePipeline(candleChan InCandles, stages ...Stage) (OutCandles, OutErrors) {
	out := make(chan dom.Candle)
	err := make(chan error)

	go func() {
		defer close(out)
		defer close(err)

		// setup pipeline
		stageOut := candleChan
		for _, stage := range stages {
			stageOut = stage(stageOut, err)
		}

		for candle := range stageOut {
			out <- candle
		}
	}()

	return out, err
}
