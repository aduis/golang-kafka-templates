Consumer Cluster Wordcount
=====
Consumer Groups are we handle scale out of partitions on the consumer.  

Consumer Group
- Consumers can join a group by using the samegroup.id.
- The maximum parallelism of a group is that the number of consumers in the group ← no of partitions.
- Kafka assigns the partitions of a topic to the consumer in a group, so that each partition is consumed by exactly one consumer in the group.
- Kafka guarantees that a message is only ever read by a single consumer in the group.
- Consumers can see the message in the order they were stored in the log.

Re-balancing of a Consumer
Adding more processes/threads will cause Kafka to re-balance. If any consumer or broker fails to send heartbeat to ZooKeeper, then it can be re-configured via the Kafka cluster. During this re-balance, Kafka will assign available partitions to the available threads, possibly moving a partition to another process.

Consumer Groups are we handle scale out of partitions on the consumer.
- Logical name for 1 or more consumers
- Message consumption is load balanced across all consumers in a group.
- Keeps only a node to process data within a consumer group.  


![consumer-groups]https://cdn.rawgit.com/wadearnold/golang-kafka-templates/9671b7eb/consumer-cluster-wordcount/image/consumer-groups.svg)
<!--
Image created with mermaid
https://knsv.github.io/mermaid/#mermaid
 source: ./image/consumer-groups.txt
-->

In this example we have two servers called Kafka brokers. The Kafka cluster has four partitions of a topic. Remember that Kafka can rebalance where those partitions are located across brokers and it automatically rebalances them. We have two logical applications that are represented by Consumer Group A and Consumer Group B. The logical application of Consumer Group A has two actual application deployable that are configured exactly the same. They could be on different physical nodes but by having the same consumer group the load of the Kafka topic is distributed to each application. In the case of word count C1 will receive words from partition 0 and 3 and C2 will receive words from partition 1 and 2. In the case that C1 were to crash C2 would be rebalanced and receive messages from partition 0-3. If another application was added into Consumer Group A for a total of three partitions the load would be rebalanced across the partitions. In this example a totality of 4 applications, one for each partition, can be added to balance the load. It is important to note that if a fifth application was added to the consumer group it would not receive any messages but would be elected to receive messages if one of the other applications in the consumer group stop receiving a heartbeat.

We have two consumer groups in this example to show that Consumer Group A is designed to be fault tolerant. Meaning any application in Group A can handle the complete load of the Kafka topic. In Consumer Group B the application requires more processing power for the same number of messages. Adding four applications to handle this work load is required to have the computational capacity to handle all the messages.     
