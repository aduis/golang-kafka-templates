Consumer Word Count
=====


This consumer example connects to all Kafka partitions and listens for messages. All messages are broken into words and a count of each word is calculated.

The mental model that I use for this client is a fanout with each client receiving all broadcasted messages on the topic.

![fanout](images/fanout.svg?raw=true)

<!--
Image created with mermaid
https://knsv.github.io/mermaid/#mermaid
 source: ./image/fanout.txt
-->


Each consumer client will have received all of the messages sent to any of the partitions available for the topic.  

I find this pattern useful for integrations. For example listening for a specific message and notify external systems such as email, SMS, slack, etc. Each of these integrations can work independently and receive all of the messages from the topic. If a services needs to be restarted it just needs to know the offset when it was shutdown. On restart it can process any messages missed while it was down. This is not a consumer client pattern for high availability or load balancing.

When running this example you will need to add messages onto the Topic. The easiest way
to do this is by running the Kafka Console Producer in a independent terminal.

```bash
localhost:~ me$ kafka-console-producer --broker-list localhost:9092 --topic test
```

Simply type words into the consuls topic and watch them get counted.  

Example output:

```bash
-- Word count current metrics --
a: 3
asdf: 3
hello: 5
world: 5
message: 3
wade: 4
```
