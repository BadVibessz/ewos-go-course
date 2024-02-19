package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	priceproducer "messageBrokers4/internal/infrastructure/producer"
)

var brokers = []string{":9093"}

func main() {
	logger := log.New()

	producer, err := priceproducer.NewPricesProducer(brokers, logger)
	if err != nil {
		logger.Fatalf("producer not started: %v", err)
	}

	ctx := context.Background() // todo: graceful shutdown

	producer.ProducePrices(ctx, "prices")
}
