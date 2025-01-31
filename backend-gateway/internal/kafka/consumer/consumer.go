package consumer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

const op = "[consumer-group]"

type ConsumerGroup struct {
	sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topics  []string
}

func (c *ConsumerGroup) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("%s run", op)

		for {
			if err := c.ConsumerGroup.Consume(ctx, c.topics, c.handler); err != nil {
				log.Printf("%s Error from consume: %v\n", op, err)
			}
			if ctx.Err() != nil {
				log.Printf("%s ctx closed: %s\n", op, ctx.Err().Error())
				return
			}
		}
	}()
}

func NewConsumerGroup(
	brokers []string,
	groupID string,
	topics []string,
	consumerGroupHandler sarama.ConsumerGroupHandler,
) (*ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.ResetInvalidOffsets = true
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	config.Consumer.Group.Session.Timeout = 60 * time.Second
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	config.Consumer.Return.Errors = true

	// config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	cg, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		ConsumerGroup: cg,
		handler:       consumerGroupHandler,
		topics:        topics,
	}, nil
}
