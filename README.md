# **Log Ingestor and Query Interface**

### **Log Ingestor**
A microservice that can handle multiple log events at the same time. It uses Kafka for message queuing and MongoDB for storing log data. This microservice can easily scale within Kubernetes by deploy multiple Kafka producer and Kafka consumer to different worker nodes in k8s cluster. And the Kafka cluster itself can be scaled by increasing the kafka brokers and the number of replication of the topic *(See the architecture for more details)*.

### **Query Interface**
A simple web interface that allows filtering log events.

### **Architecture**
![alt text](https://github.com/phonghaido/log-ingestor/blob/main/images/artifacture.png?raw=true)
### **Usage**
##### 1.  Run docker compose file
```shell
docker-compose up -d
```

##### 2. Verify kafka and zookeeper
```shell
nc -zv localhost 9092
nc -zv localhost 2181
```
##### 3. Ingest log into system
The script will generate multiple log requests and send those requests to kafka producer concurrently. The number of request is configurable, the current number of concurrent request is 1000
```shell
go run scripts/main.go
```
##### 4. Double check the data in MongoDB with MongoDB Compass or MongoDB Shell (Optional)
MongoDB Compass
![alt text](https://github.com/phonghaido/log-ingestor/blob/main/images/db.png?raw=true)

MongoDB Shell
 ```shell
use log_ingestor
db.logs.count()
```
##### **5.  Search for logs**
![alt text](https://github.com/phonghaido/log-ingestor/blob/main/images/ui.png?raw=true)
