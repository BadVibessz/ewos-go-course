package producer

import (
	"context"
	"fmt"
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
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

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
		//priceJSON, err := json.Marshal(price)
		//if err != nil {
		//	p.log.Errorf("error occured marshalling price: %v", err)
		//}

		msg := sarama.ProducerMessage{
			Topic:    topic,
			Key:      sarama.StringEncoder("key1"),
			Value:    sarama.StringEncoder(fmt.Sprintf("new price for %s", price.Ticker)),
			Metadata: price,
		}

		partition, offset, err := p.producer.SendMessage(&msg)
		if err != nil {
			p.log.Errorf("error occured producing message: %v", err)
		}

		p.log.Infof("message produced: partition=%v, offset=%v, value: %v", partition, offset, msg.Value)

		msg = sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder("key2"),
			Value: sarama.StringEncoder("key2 msg"),
		}

		partition, offset, err = p.producer.SendMessage(&msg)
		if err != nil {
			p.log.Errorf("error occured producing message: %v", err)
		}

		p.log.Infof("message produced: partition=%v, offset=%v, value: %v", partition, offset, msg.Value)
	}

	p.log.Info("prices generator stopped")
}

func (p *Prices) Close() error {
	return p.producer.Close()
}
