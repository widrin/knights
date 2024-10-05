package clusterConsumer

import (
	"context"
	"github.com/IBM/sarama"
	"strings"
	"sync"
	"tzgit.kaixinxiyou.com/utils/common/log"
)

const ConsumerMsgEvent = "ConsumerMsgEvent"

type ClusterConsumer struct {
	brokers           []string
	topics            []string
	groupId           string
	consumer          sarama.ConsumerGroup
	closeSig          chan bool
	partitionConsumer []sarama.PartitionConsumer
	consumerData      *Consumer
	config            *sarama.Config
}

func NewClusterConsumer(brokers string, topics string, groupId string, user string, pass string) *ClusterConsumer {
	c := &ClusterConsumer{}

	c.brokers = strings.Split(brokers, ",")
	c.topics = strings.Split(topics, ",")
	c.groupId = groupId
	c.closeSig = make(chan bool, 1)
	c.consumerData = &Consumer{
		ready: make(chan bool),
	}
	c.config = sarama.NewConfig()
	if user != "" {
		c.config.Net.SASL.Mechanism = "PLAIN"
		c.config.Net.SASL.Version = int16(1)
		c.config.Net.SASL.Enable = true
		c.config.Net.SASL.User = user
		c.config.Net.SASL.Password = pass
	}

	c.config.Consumer.Offsets.Initial = sarama.OffsetNewest
	c.config.Version = sarama.V1_1_1_0
	return c
}
func NewClusterConsumer1(brokers string, topics string, groupId string, config *sarama.Config) *ClusterConsumer {
	c := &ClusterConsumer{}

	c.brokers = strings.Split(brokers, ",")
	c.topics = strings.Split(topics, ",")
	c.groupId = groupId
	c.closeSig = make(chan bool, 1)
	c.consumerData = &Consumer{
		ready: make(chan bool),
	}
	c.config = config
	return c
}
func (c *ClusterConsumer) Init(autoCommit bool) error {
	var err error

	if !autoCommit {
		c.config.Consumer.Offsets.AutoCommit.Enable = false
	}
	c.consumerData.autoCommit = autoCommit

	c.consumer, err = sarama.NewConsumerGroup(c.brokers, c.groupId, c.config)
	return err
}
func (c *ClusterConsumer) SetOnMessage(f func(*sarama.ConsumerMessage)) {
	c.consumerData.onConsumer = f
}
func (c *ClusterConsumer) Run() {

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := c.consumer.Consume(ctx, c.topics, c.consumerData); err != nil {
				log.Warn("Error from consumer: %v", err)
				break
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			c.consumerData.ready = make(chan bool)
		}
	}()
	<-c.consumerData.ready // Await till the consumer has been set up
	select {
	case <-ctx.Done():
		log.Debug("terminating: context cancelled")
	case <-c.closeSig:
		log.Debug("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err := c.consumer.Close(); err != nil {
		log.Warn("Error closing client: %v", err)
	}
}

func (c *ClusterConsumer) Close() {
	c.closeSig <- true
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready      chan bool
	onConsumer func(*sarama.ConsumerMessage)
	autoCommit bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if consumer.onConsumer != nil {
			consumer.onConsumer(message)
		}
		//skeleton.OnChanRpcEvent(ConsumerMsgEvent, message)
		session.MarkMessage(message, "")
	}

	if !consumer.autoCommit {
		session.Commit()
	}

	return nil
}
