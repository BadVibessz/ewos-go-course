package pipeline

import (
	"context"
	dom "hw-async/domain"
)

type (
	InCandles = <-chan dom.Candle
	InErrors  = chan error

	OutCandles = InCandles
	OutErrors  = InErrors
)

type Stage func(InCandles, InErrors) OutCandles

func ExecutePipeline(ctx context.Context, candleChan InCandles, stages ...Stage) (OutCandles, OutErrors) {
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

		for {
			select {
			case val, open := <-stageOut: // todo: maybe use for range?
				if open {
					out <- val
				} else {
					return
				}

			case <-ctx.Done():
				//val, open := <-stageOut
				//fmt.Printf("%+v\n, %s", val, open)
				//
				//if open {
				//	out <- val
				//} else {
				//	return
				//}

				// todo: handle graceful shutdown
				// return
			}
		}
	}()

	return out, err
}
