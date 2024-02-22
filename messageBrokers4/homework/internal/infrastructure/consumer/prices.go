package consumer

import (
	"github.com/IBM/sarama"
	"time"

	log "github.com/sirupsen/logrus"
)

type Prices struct {
	ready  chan bool
	logger *log.Logger
}

func NewPricesConsumer(logger *log.Logger) *Prices {
	return &Prices{
		ready:  make(chan bool),
		logger: logger,
	}
}

func (p *Prices) Setup(session sarama.ConsumerGroupSession) error {
	close(p.ready) // todo:
	return nil
}

func (p *Prices) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil // todo:
}

func (p *Prices) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				p.logger.Infof("message channel was closed")

				return nil
			}

			p.logger.Infof("Message claimed: value = %s, timestamp = %v, topic = %s partition: %v",
				string(message.Value), message.Timestamp.Format(time.ANSIC), message.Topic, message.Partition)

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			p.logger.Infof("consume claim: parent context was cancelled")

			return nil
		}
	}
}
