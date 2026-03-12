package publisher

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type KafkaPublisher struct {
	producer sarama.SyncProducer
}

func NewKafkaPublisher(brokers []string) (*KafkaPublisher, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &KafkaPublisher{
		producer: producer,
	}, nil
}

func (k *KafkaPublisher) Publish(topic, key string, message []byte, headers map[string]string) error {
	saramaHeaders := make([]sarama.RecordHeader, 0, len(headers))
	for k, v := range headers {
		saramaHeaders = append(saramaHeaders, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}
	msg := &sarama.ProducerMessage{
		Topic:   topic,
		Key:     sarama.StringEncoder(key),
		Value:   sarama.ByteEncoder(message),
		Headers: saramaHeaders,
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to publish message to Kafka: %w", err)
	}

	log.Printf("Message published to Kafka: topic=%s, partition=%d, offset=%d, key=%s",
		topic, partition, offset, key)

	return nil
}

func (p *KafkaPublisher) Close() error {
	if err := p.producer.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka producer: %w", err)
	}
	log.Println("Kafka producer closed")
	return nil
}
