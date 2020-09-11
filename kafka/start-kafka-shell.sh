#!/bin/bash
# $1 = [隨便給個ip] 
# $2 = [給zoo1:2181] 
docker run --rm --network jay-play-go_proxy -v /var/run/docker.sock:/var/run/docker.sock -e HOST_IP=$1 -e ZK=$2 -i -t wurstmeister/kafka /bin/bash

### 建立 topic ###
$KAFKA_HOME/bin/kafka-topics.sh --create --topic topic \
--partitions 4 --zookeeper $ZK --replication-factor 2

### 看看現在有哪些 topic ###
$KAFKA_HOME/bin/kafka-topics.sh --list  --zookeeper $ZK

### 看看現在有哪些 topic 的資訊 ###
$KAFKA_HOME/bin/kafka-topics.sh --describe --topic topic --zookeeper $ZK

### 發送訊息 ###
$KAFKA_HOME/bin/kafka-console-producer.sh --topic=topic --broker-list=`broker-list.sh`

### 接收訊息 ###
$KAFKA_HOME/bin/kafka-console-consumer.sh --topic=topic --bootstrap-server=<broker location> --from-beginning
