package main

import (
	"log"
	"notification-service/controllers"
	"notification-service/database"
	"notification-service/utils"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()
	db := database.GetDB()

	// Initialize Redis
	utils.InitRedis()

	// Initialize Kafka producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": ":9092",
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// starting kafka
	// kafka_2.13-3.9.0/bin/zookeeper-server-start.sh kafka_2.13-3.9.0/config/zookeeper.properties
	// kafka_2.13-3.9.0/bin/kafka-server-start.sh kafka_2.13-3.9.0/config/server.properties

	// Define Kafka topic
	topic := "notification.send_sms"

	utils.InitElasticsearch()
	// Initialize Gin router
	router := gin.Default()

	// Initialize controllers
	smsController := &controllers.SMSController{
		Producer: producer,
		DB:       db,
		Topic:    topic,
	}

	blacklistController := &controllers.BlacklistController{
		DB: db,
	}

	// Define routes
	router.POST("/v1/sms/send", smsController.SendSMSHandler)
	router.GET("/v1/sms/:requestId", smsController.GetSMSDetails)
	router.GET("/v1/blacklist", blacklistController.GetBlacklistedNumbers)
	router.POST("/v1/blacklist", blacklistController.AddToBlacklist)
	router.DELETE("/v1/blacklist", blacklistController.RemoveFromBlacklist)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		_, err := utils.RedisClient.Ping(utils.Ctx).Result()
		if err != nil {
			c.JSON(500, gin.H{"status": "DOWN", "error": "Redis not connected"})
			return
		}
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Start the HTTP server
	log.Println("Starting Notification Service on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
