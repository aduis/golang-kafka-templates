Consumer Word Count
=====

This consumer examples connects to all Kafka partitions and listens for messages.
All messages are broken into words and a count of each word assurance is calculated.

When running this example you will need to add messages onto the Topic. The easiest way
to do this is by running the Kafka Console Producer in a independent terminal.

```bash
kafka-console-producer --broker-list localhost:9092 --topic test
```

Simply type words into the topic and watch them get counted.  
