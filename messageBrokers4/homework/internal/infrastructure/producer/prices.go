package producer

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/IBM/sarama"

	"messageBrokers4/internal/domain/generator"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

type Prices struct {
	log      *log.Logger
	producer sarama.SyncProducer
}

func NewPricesProducer(brokers []string, log *log.Logger) (*Prices, error) {
	config := sarama.NewConfig()
	// config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	res := Prices{
		log:      log,
		producer: producer,
	}

	return &res, nil
}

func (p *Prices) ProducePrices(ctx context.Context, topic string) {
	pg := generator.NewPrices(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})

	p.log.Info("start prices generator...")
	prices := pg.GeneratePrices(ctx)

	for price := range prices {
		priceJSON, err := json.Marshal(price)
		if err != nil {
			p.log.Errorf("error occured marshalling price: %v", err)
		}

		msg := sarama.ProducerMessage{
			Topic:     topic,
			Key:       nil,
			Value:     sarama.StringEncoder(priceJSON),
			Offset:    0,
			Partition: 0,
		}

		partition, offset, err := p.producer.SendMessage(&msg)
		if err != nil {
			p.log.Errorf("error occured producing message: %v", err)
		}

		p.log.Infof("message produced: partition=%v, offset=%v", partition, offset)
	}

	p.log.Info("prices generator stopped")
}
