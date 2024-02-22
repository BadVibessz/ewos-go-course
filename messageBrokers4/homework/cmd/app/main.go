package main

import (
	"context"
	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
	"messageBrokers4/internal/infrastructure/consumer"
	"messageBrokers4/internal/infrastructure/producer"
)

var (
	brokers = []string{":9093"}
	topics  = []string{"prices"}
)

const consumerGroupName = "my_group"

func initPriceProducer(logger *log.Logger) *producer.Prices {
	prod, err := producer.NewPricesProducer(brokers, logger)
	if err != nil {
		logger.Fatalf("prod not started: %v", err)
	}

	return prod
}

func initPriceConsumer() {

}

func main() {
	logger := log.New()
	ctx := context.Background() // todo: graceful shutdown

	// price producer
	prod := initPriceProducer(logger)

	defer func() {
		if err := prod.Close(); err != nil {
			log.Fatalf("failed to close producer")
		}
	}()

	prod.ProducePrices(ctx, "prices")

	// price consumer
	conf := sarama.NewConfig()

	consGroup, err := sarama.NewConsumerGroup(brokers, consumerGroupName, conf)
	if err != nil {
		logger.Fatalf("cannot create consumer group: %v", err)
	}

	// track errors
	go func() {
		for err = range consGroup.Errors() {
			logger.Errorf("error occurred in consumer group %s: %v", consumerGroupName, err)
		}
	}()

	// iterate over consumer sessions.
	logger.Infof("price consumer started")
	for {
		handler := consumer.NewPricesConsumer(logger)

		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		err = consGroup.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}

}
