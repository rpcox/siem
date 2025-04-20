## siem : kafka

This is currently running on three Raspberry Pis, Pi 4 Model B w/ 4 GB RAM and 128 GB Extreme Sandisk Micro SD cards.

The components are:

* Debian GNU/Linux 12 (bookworm)
* OpenJDK 17.0.14
* Kafka 2.13-4.0.0

![image](png/dragons.drawio.png)


### installation

1. Install OpenJDK
2. Install Kakfka
3. Install systemd unit file and leave disabled
4. Configure the broker/controllers (see configs)
* Comment out advertised.listeners
6. Generate a cluster ID
* KAFKA_CLUSTER_ID=$(bin/kafka-storage.sh random-uuid)
* Keep in case you need to reformat one node
6. Format storage
* bin/kafka-storage.sh format -t $KAFKA_CLUSTER_ID -c config/server.properties
7. Start Kafka on each node


### test

All action from puff


Create a topic

    > bin/kafka-topics.sh --create --topic test.topic --bootstrap-server mushu:9092
    Created topic test.topic
    >


List the new topic

    > bin/kafka-topics.sh --bootstrap-server mushu:9092 --list --exclude-internal
    test.topic
    >


Pull the topic description

    > bin/kafka-topics.sh --describe --bootstrap-server smaug:9092 --topic test.topic
    Topic: test.topic	TopicId: QmLwMsncTBWiSn2HoB0Kqg	PartitionCount: 6	ReplicationFactor: 1	Configs: segment.bytes=1073741824
	Topic: test.topic	Partition: 0	Leader: 300	Replicas: 300	Isr: 300	Elr: 	LastKnownElr:
	Topic: test.topic	Partition: 1	Leader: 100	Replicas: 100	Isr: 100	Elr: 	LastKnownElr:
	Topic: test.topic	Partition: 2	Leader: 200	Replicas: 200	Isr: 200	Elr: 	LastKnownElr:
	Topic: test.topic	Partition: 3	Leader: 300	Replicas: 300	Isr: 300	Elr: 	LastKnownElr:
	Topic: test.topic	Partition: 4	Leader: 200	Replicas: 200	Isr: 200	Elr: 	LastKnownElr:
	Topic: test.topic	Partition: 5	Leader: 100	Replicas: 100	Isr: 100	Elr: 	LastKnownElr:
    >


Produce

    > kafka-console-producer.sh --topic test.topic --bootstrap-server mushu:9092
    >one
    >two
    >three
    >^C


Consume

    > kafka-console-consumer.sh --topic test.topic --bootstrap-server smaug:9092 --from-beginninging
    one
    two
    three
    ^CProcessed a total of 3 messages
    >


Delete the topic

    > bin/kafka-topics.sh --bootstrap-server mushu:9092 --delete --topic test.topic
    >


Verify deletion

     > kafka-topics.sh --bootstrap-server mushu:9092 --list --exclude-internal
     >





