package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

// serviceName of the service being run. Appended to all logs.
const serviceName = "wordcount-api"

// topic is the Kafka topic to connect to for messags
const topic = "wordcount-topic"

func main() {

	// Get environment variables.
	var brokerList = []string{"localhost:9092"}
	if os.Getenv("brokerList") != "" {
		brokerList = []string{os.Getenv("brokerList")}
	}

	// create a syslog
	logger, err := syslog.New(syslog.LOG_LOCAL3, serviceName)
	if err != nil {
		panic("Cannot attach to syslog: " + err.Error())
	}
	defer logger.Close()
	logger.Warning("Starting " + serviceName + " producer")

	// build the Sarama Configuration
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		logger.Warning("Unable to create a Kafka producer client " + err.Error())
		panic("Unable to create a Kafka producer Client " + err.Error())
	}

	// start the http services

	// landing page get request
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(w, "HTTP endpoint for Kafka ")
	})

	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		r.ParseForm()
		// so much logic needs to be here to ENSURE you want this message. Consider it a
		// database write.
		msg := buildMessage(r.FormValue("q"), topic)
		_, _, err = producer.SendMessage(msg)
		if err != nil {
			logger.Warning("Could not send message: " + err.Error())
		}
	})

	// start the http server
	log.Fatal(http.ListenAndServe(":8081", nil))

	// handle graceful shutdown
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)

	for {
		select {
		case <-signals:
			logger.Warning("Interrupt shutdown is detected")
			producer.Close()
			return
		}
	}

}

func buildMessage(msg string, topic string) *sarama.ProducerMessage {
	// This section can be expand to tweak the producer client message
	message := &sarama.ProducerMessage{
		Topic: topic,
		// Partition: partition,
		Value: sarama.StringEncoder(msg),
	}
	return message
}
