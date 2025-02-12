# Notification Service

This is a Go-based microservice responsible for handling notifications. It consumes messages from Kafka, processes them, and integrates with a third-party SMS provider to send notifications.

## Features
- Kafka consumer for handling notification requests
- Redis caching for blacklisted numbers
- SQL database for storing SMS requests
- Elasticsearch integration for message indexing
- 3rd-party API integration for SMS delivery

## Prerequisites
Ensure you have the following installed on your system:
- **Go** (>= 1.18)
- **Docker & Docker Compose**
- **Kafka** (Running on localhost:9092)
- **Redis** (Running on localhost:6379)
- **SQL Database** (Running on localhost)
- **Elasticsearch** (Running on localhost:9200)

## Installation & Setup

### 1️⃣ Clone the Repository
```sh
git clone git@github.com:Latanshu-meesho/notification-service.git
cd notification-service
```

### 2️⃣ Set Up Environment Variables
Create a `.env` file in the project root and configure the necessary values:
```env
KAFKA_BROKER=localhost:9092
REDIS_ADDR=localhost:6379
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=notification_service
ELASTICSEARCH_URL=http://localhost:9200
SMS_API_KEY=your_sms_provider_key
```

### 3️⃣ Start Dependencies (Using Docker Compose)
```sh
docker-compose up
```

### 4️⃣ Install Dependencies
```sh
go mod tidy
```

### 5️⃣ Run the Notification Service
```sh
go run cmd/main.go
```

## Commands for Setup

### Start Kafka, Redis, SQL, and Elasticsearch (Docker)
```sh
docker-compose up -d
```

### Create Kafka Topic (If not already created)
```sh
kafka-topics --create --topic notification.send_sms --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
```

### Check Kafka Messages
```sh
kafka-console-consumer --bootstrap-server localhost:9092 --topic notification.send_sms --from-beginning
```

### Verify Redis Connection
```sh
redis-cli
ping
```

### Verify SQL Connection
```sh
mysql -h localhost -u your_username -p -D notification_service
```

### Check Elasticsearch Health
```sh
curl -X GET "http://localhost:9200/_cluster/health?pretty"
```

## Usage
### 📩 Sending an SMS Request
The service listens to the Kafka topic `notification.send_sms`. To test, produce a message:
```sh
kafka-console-producer --broker-list localhost:9092 --topic notification.send_sms
```
Then, enter a message containing a `requestID`.

### 🔍 Searching SMS Logs in Elasticsearch
```sh
curl -X GET "http://localhost:9200/sms_logs/_search" -H "Content-Type: application/json" -d '{
  "query": {
    "match": { "text": "your_message_content" }
  }
}'
```

## Troubleshooting
- **Kafka Consumer Not Reading Messages?**
  - Ensure Kafka is running: `docker ps | grep kafka`
  - Check consumer group status: `kafka-consumer-groups --bootstrap-server localhost:9092 --group notification-consumer-group --describe`

- **Redis Not Storing Blacklisted Numbers?**
  - Run `redis-cli` and check keys: `KEYS *`

- **SQL Database Connection Issues?**
  - Ensure database is running and credentials are correct.

## Contributing
Feel free to open issues and submit pull requests.