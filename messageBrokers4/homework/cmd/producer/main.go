package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"messageBrokers4/internal/infrastructure/producer"
)

var brokers = []string{":9093"}

func main() {
	logger := log.New()

	prod, err := producer.NewPricesProducer(brokers, logger)
	if err != nil {
		logger.Fatalf("prod not started: %v", err)
	}

	ctx := context.Background() // todo: graceful shutdown

	prod.ProducePrices(ctx, "prices")

	defer func() {
		if err = prod.Close(); err != nil {
			log.Fatalf("failed to close producer")
		}
	}()
}
