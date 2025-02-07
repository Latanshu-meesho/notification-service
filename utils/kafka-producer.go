package utils

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func PublishKafkaMessage(producer *kafka.Producer, topic string, value string) error {
	deliveryChan := make(chan kafka.Event)

	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
	}, deliveryChan)

	if err != nil {
		log.Printf("Failed to produce Kafka message: %v", err)
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Failed to publish to kafkaTopic: %v", m.TopicPartition.Error)
	} else {
		log.Printf("RequestID : %s published to topic %s [%d] at offset %v", m.Value, *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
	return nil
}
