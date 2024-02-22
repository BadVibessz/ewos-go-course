package main

import (
	"context"
	"github.com/IBM/sarama"
	"messageBrokers4/internal/infrastructure/consumer"

	log "github.com/sirupsen/logrus"
)

var brokers = []string{":9093"}

const consumerGroupName = "my_group"

func main() {
	logger := log.New()

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
	ctx := context.Background()

	// TODO: graceful shutdown

	logger.Infof("price consumer started")
	for {
		topics := []string{"prices"}
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
