package main

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/memphisdev/superstream.go"
)

func main() {
	broker := "..."
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Flush.MaxMessages = 10
	config.Producer.RequiredAcks = sarama.NoResponse

	// confluent config
	config.Net.SASL.Enable = true
	config.Net.SASL.User = "..."
	config.Net.SASL.Password = "..."
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = nil

	config = superstream.Init("token", config, superstream.ConsumerGroup("group"), superstream.Servers(broker))

	consumer, err := sarama.NewConsumerGroup([]string{broker}, "group", config)
	if err != nil {
		panic(err)
	}

	kafkaHandler := KafkaConsumerGroupHandler{}

	for {
		err := consumer.Consume(context.Background(), []string{"test"}, &kafkaHandler)
		if err != nil {
			panic(err)
		}
	}
}

type KafkaConsumerGroupHandler struct{}

func (h *KafkaConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *KafkaConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *KafkaConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg := <-claim.Messages():
			fmt.Print(string(msg.Value))
			sess.MarkMessage(msg, "")

		case <-sess.Context().Done():
			return nil
		}
	}
}