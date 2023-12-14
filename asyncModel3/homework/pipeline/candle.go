package pipeline

import (
	"context"
	dom "hw-async/domain"
	_ "hw-async/generator"
)

type (
	InPrice  = <-chan dom.Price
	InCandle = <-chan dom.Candle
	Out      = InCandle
)

type Stage func(context.Context, InPrice, InCandle) Out

func ExecutePipeline(ctx context.Context, priceChan InPrice, stages ...Stage) Out {
	out := make(chan dom.Candle)

	go func() {
		defer close(out)

		// execute pipeline
		stageOut := make(InCandle)
		for _, stage := range stages {
			stageOut = stage(ctx, priceChan, stageOut)
		}

		for {
			select {
			case val, open := <-stageOut:
				if open {
					out <- val
				} else {
					return
				}

			case <-ctx.Done():
				// todo: handle graceful shutdown
				return
			}
		}
	}()

	return out
}
