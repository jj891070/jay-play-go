### 檢查zookeeper是否健康 ###
$ docker-compose logs zookeeper | grep -i binding

### 檢查kafka是否健康 ###
$ docker-compose logs kafka | grep -i started

### 建立topic ###
$ docker-compose exec kafka  \
kafka-topics --create --topic foo --partitions 1 --replication-factor 1 --if-not-exists --zookeeper zookeeper:2181

### 檢查建立topic ###
$ docker-compose exec kafka  \
  kafka-topics --describe --topic foo --zookeeper localhost:32181

### 發送訊息 ###
$ docker-compose exec kafka  \
  bash -c "seq 42 | kafka-console-producer --request-required-acks 1 --broker-list localhost:9092 --topic foo && echo 'Produced 42 messages.'"

### 接收訊息 ###
$ docker-compose exec kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic foo --from-beginning --max-messages 42