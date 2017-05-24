package main

import (
	"fmt"

	"github.com/Shopify/sarama"
)

const topic = "test"

var brokers = []string{"127.0.0.1:9092"}

func main() {
	fmt.Println("A Kafka consumer example")
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		fmt.Println("Could not create consumer: ", err)
	}
	subscribe(topic, consumer)
	select {}
}

func subscribe(topic string, consumer sarama.Consumer) {
	partionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println("Error retrieving partition list: ", err)
	}
	initialOffset := sarama.OffsetOldest
	for _, paritition := range partionList {
		pc, _ := consumer.ConsumePartition(topic, paritition, initialOffset)
		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				messageRecieved(message)
			}
		}(pc)
	}
}

func messageRecieved(message *sarama.ConsumerMessage) {
	fmt.Printf("Message: %v \n", string(message.Value))
}
