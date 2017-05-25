package main

import (
	"fmt"
	"log/syslog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

// serviceName of the service being run. Appended to all logs.
const serviceName = "wordcounter"

// topic is the Kafka topic to connect to for messags
const topic = "wordcount-topic"

func main() {
	// Get environment variables.
	var brokerList = []string{"localhost:9092"}
	if os.Getenv("brokerList") != "" {
		brokerList = []string{os.Getenv("brokerList")}
	}
	// offset -1 = OffsetNewest | -2 = OffsetOldest
	var offset int64 = -2
	if os.Getenv("offset") != "" {
		offset, _ = strconv.ParseInt(os.Getenv("offset"), 10, 64)
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

	// Add shared libraries for each goroutine. This is buisness logic for the word count example
	w := newWordCountMap()
	// buisness logic requires printing current counts on an intraval
	ticker := time.NewTicker(30 * time.Second)

	for _, paritition := range partitionList {
		pc, _ := consumer.ConsumePartition(topic, paritition, offset)

		go func(pc sarama.PartitionConsumer, w *wordCountMap, logger *syslog.Writer) {
			for message := range pc.Messages() {
				messageRecieved(message, w) //  could also use a channel
			}
			for {
				select {
				case err := <-pc.Errors():
					logger.Warning("Partition error: " + err.Error())
				}
			}
		}(pc, w, logger)
	}

	for {
		select {
		case <-signals:
			logger.Warning("Interrupt shutdown is detected")
			consumer.Close()
			ticker.Stop()
			return
		case <-ticker.C:
			fmt.Println("-- Word count current metrics --")
			w.Lock()
			for word, count := range w.found {
				if count > 1 {
					fmt.Printf("%s: %d\n", word, count)
				}
			}
			w.Unlock()
			logger.Warning("Word Count was calculated")
		}
	}
}

// messageRecieved is where you add the buisness logic for what to do with messages that come
// from the topic that the consumer is listening on.
func messageRecieved(message *sarama.ConsumerMessage, w *wordCountMap) {
	fmt.Printf("Message: %v \n", string(message.Value))
	// arbitray buisness logic
	words := strings.Fields(string(message.Value))
	for _, word := range words {
		w.add(strings.ToLower(word), 1)
	}
}

// wordCountMap is a shared map of all words that have been seen.
type wordCountMap struct {
	sync.Mutex
	found map[string]int
}

// newWordCountMap returns an instantiated wordCountMap
func newWordCountMap() *wordCountMap {
	return &wordCountMap{found: map[string]int{}}
}

// add appends words into the wordCountMap
func (w *wordCountMap) add(word string, n int) {
	w.Lock()
	defer w.Unlock()
	count, ok := w.found[word]
	if !ok {
		w.found[word] = n
		return
	}
	w.found[word] = count + n
}
