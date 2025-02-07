package utils

import (
	"log"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     ":6379", // Replace with your Redis server address
		Password: "",      // Replace with your Redis password
		DB:       0,       // Default DB
	})

	// Test the connection
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}

func IsBlacklisted(phoneNumber string) (bool, error) {
	result, err := RedisClient.SIsMember(Ctx, "blacklist", phoneNumber).Result()
	if err != nil {
		log.Printf("Failed to check blacklist: %v", err)
		return false, err
	}
	return result, nil
}

func AddToRedisBlacklist(phoneNumbers []string) error {
	for _, number := range phoneNumbers {
		if err := RedisClient.SAdd(Ctx, "blacklist", number).Err(); err != nil {
			log.Printf("Failed to add number to Redis blacklist: %v", err)
			return err
		}
	}
	return nil
}

func RemoveFromRedisBlacklist(phoneNumbers []string) error {
	for _, number := range phoneNumbers {
		if err := RedisClient.SRem(Ctx, "blacklist", number).Err(); err != nil {
			log.Printf("Failed to remove number from Redis blacklist: %v", err)
			return err
		}
	}
	return nil
}

func GetBlacklistFromRedis() ([]string, error) {
	result, err := RedisClient.SMembers(Ctx, "blacklist").Result()
	if err != nil {
		log.Printf("Failed to get blacklist from Redis: %v", err)
		return nil, err
	}
	return result, nil
}
