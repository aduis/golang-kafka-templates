package main

import (
	"fmt"
	"log/syslog"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

// serviceName of the service being run. Appended to all logs.
const serviceName = "consumer-min"

// topic is the Kafka topic to connect to for messags
const topic = "test"

func main() {
	// Get environment variables.
	var brokerList = []string{"localhost:9092"}
	if os.Getenv("brokerList") != "" {
		brokerList = []string{os.Getenv("brokerList")}
	}
	// offset -1 = OffsetNewest | -2 = OffsetOldest
	var offset int64 = -2
	if os.Getenv("offset") != "" {
		// offset = strconv.Atoi(os.Getenv("offset"))
	}

	// create a syslog
	logger, err := syslog.New(syslog.LOG_LOCAL3, serviceName)
	if err != nil {
		panic("Cannot attach to syslog: " + err.Error())
	}
	defer logger.Close()
	logger.Warning("Starting " + serviceName + " consumer")

	// handle graceful shutdown
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)
	//messages := make(chan *sarama.ConsumerMessage, 256)

	// configure the kafka connection
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		logger.Warning("Could not create consumer: " + err.Error())
		panic("Could not create consumer: " + err.Error())
	}
	// partitionList returns the sorted list of all partition ID's for the given topic
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		logger.Warning("Error retrieving partition list: " + err.Error())
		panic("Error retrieving partition list: " + err.Error())
	}

	for _, paritition := range partitionList {
		pc, _ := consumer.ConsumePartition(topic, paritition, offset)

		go func(pc sarama.PartitionConsumer, logger *syslog.Writer) {
			for message := range pc.Messages() {
				messageRecieved(message) //  could also use a channel
			}
			for {
				select {
				case err := <-pc.Errors():
					logger.Warning("Partition error: " + err.Error())
				}
			}
		}(pc, logger)
	}

	for {
		select {
		case <-signals:
			logger.Warning("Interrupt is detected")
			consumer.Close()
			return
		}
	}
}

// messageRecieved is where you add the buisness logic for what to do with messages that come
// from the topic that the consumer is listening on.
func messageRecieved(message *sarama.ConsumerMessage) {
	fmt.Printf("Message: %v \n", string(message.Value))
}
