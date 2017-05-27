Golang Kafka Microservice Templates
=====

A group of [GoLang](https://golang.org/) example templates for writing [Apache Kafka](https://kafka.apache.org/) production ready applications. The templates were created to be used as a starting point for [microservices](https://martinfowler.com/articles/microservices.html) that utilize Kafka for [Event Sourcing](https://martinfowler.com/eaaDev/EventSourcing.html) as an eventstore. An eventstore is the backbone of an event-sourced mircroservice architecture. Utilizing best practices in GoLang for [producers](https://kafka.apache.org/documentation/#producerapi) and [consumers](https://kafka.apache.org/documentation/#consumerapi) is therefore critical to a successful project.    

Each template does not have any domain specific logic. In some cases wordcount is implemented just to show were to place logic and how to do it in a thread safe fashion. The goal is that new developers can leverage the boiler plate code examples to add additional functionality to an application.

### Example Templates
- [consumer-wordcount](./consumer-wordcount/) is a minimal implementation of a Kafka consumer that provides:
  - Graceful shutdown to avoid data loss and unexpected behavior
  - Configuration via environment variables
  - Syslog logging of information and issues
  - Implements a word count logic for a basic accumulator of messages

- [producer-wordcount-api](./producer-wordcount-api) is a minimal implementation of a Kafka producer for sending http post vars into a a Kafka topic.
  - Graceful shutdown to avoid data loss and unexpected behavior
  - Configuration via environment variables
  - Syslog logging of information and issues

### Usage
It assumed that you have an understanding of [Apache Kafka](https://kafka.apache.org/) and have a version of Kafka running that you can connect to for development. The Kafka [quickstart](https://kafka.apache.org/quickstart) tutorial can assist you in running a local copy. All examples require Kafka v0.9+ to utilize Kafka based client tracking rather than zookeeper.    

On macOS with [Homebrew](https://brew.sh/) installed the following starts a local copy of Kafka

```bash
localhost:~ me$ Brew install Kafka
localhost:~ me$ zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties & kafka-server-start /usr/local/etc/kafka/server.properties
```

The code requires [Sarama](https://github.com/Shopify/sarama) an MIT-licensed Go client library for Apache Kafka. We utilize this library because it is a pure go implementaion of the wire transfer protocol of Kafka with no C dependencies. It also takes advantage of [consumer groups](http://kafka.apache.org/documentation.html#impl_zkconsumers) which ensures that for replicated consumers duplicated execution of published messages doesn't occur.

```bash
localhost:~ me$ go get github.com/Shopify/sarama
```

Clustered consumer examples utilize [Sarama-Cluster](https://github.com/bsm/sarama-cluster) a cluster Sarama extension enabling users to consume topics across balanced nodes. It's important to read the [Client-side Assignment Proposal](https://cwiki.apache.org/confluence/display/KAFKA/Kafka+Client-side+Assignment+Proposal) to understand use cases for a clustered consumer. 

```bash
localhost:~ me$ go get github.com/bsm/sarama-cluster
```



All [logs](https://golang.org/pkg/log/syslog/) are sent to [syslog](https://en.wikipedia.org/wiki/Syslog)

```bash
localhost:~ me$ tail -f /var/log/system.log | grep "serviceName"
```

Please do not fork the repository unless you are submitting a pull request to update the template. For your own usage please [mirror](https://help.github.com/articles/duplicating-a-repository/) the repo.

### Design Assumptions
These templates follow [The Twelve Factors](https://12factor.net/) a popular and widely used methodology for building web applications.


### Thank you
This work is built on top of the incredible blog posts from [Oscar Oranagwa](https://medium.com/@Oskarr3)

- [Building Scalable Applications Using Event Sourcing and CQRS](https://medium.com/technology-learning/event-sourcing-and-cqrs-a-look-at-kafka-e0c1b90d17d8)
- [Implementing CQRS using Kafka and Sarama Library in Golang](https://medium.com/@Oskarr3/implementing-cqrs-using-kafka-and-sarama-library-in-golang-da7efa3b77fe)
