package scheduler

import (
	"github.com/GGmaz/wallet-arringo/pkg/wire"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"strings"
	"time"
)

// StartDataCollector starts a background task to collect data from Redis.
func StartDataCollector(ctx *gin.Context, redisClient *redis.Client) {
	go func() {
		for {
			// Fetch accounts and transactions stored for more than 50 seconds
			fetchOldData(ctx, redisClient)

			// Sleep for 10 seconds before the next iteration
			time.Sleep(10 * time.Second)
		}
	}()
}

// fetchOldData retrieves accounts and transactions stored for more than 50 seconds from Redis.
func fetchOldData(ctx *gin.Context, redisClient *redis.Client) {
	// Fetch and process old accounts
	processOldData(ctx, "accountNum:", redisClient)

	// Fetch and process old transactions
	processOldData(ctx, "transactionId:", redisClient)
}

// processOldData fetches data from Redis based on the given key pattern and processes it if its value is older than 50 seconds.
func processOldData(ctx *gin.Context, keyPattern string, redisClient *redis.Client) {
	// Get all keys matching the given pattern
	keys, err := redisClient.Keys(ctx, keyPattern+"*").Result()
	if err != nil {
		log.Println("Error fetching keys from Redis:", err)
		return
	}

	// Iterate through keys
	for _, key := range keys {
		// Get value (time) associated with the key
		value, err := redisClient.Get(ctx, key).Int64()
		if err != nil {
			log.Println("Error fetching value from Redis for key", key, ":", err)
			continue
		}

		// Calculate age of the data (current time - stored time)
		currentTime := time.Now().Unix()
		age := currentTime - value

		// Check if data is older than 50 seconds
		if age > 50 {
			restOfString := strings.TrimPrefix(key, keyPattern)

			if keyPattern == "accountNum:" {
				err = wire.Svc.AccountService.VerifyAccount(ctx, restOfString)
				if err != nil {
					continue
				}
			} else if keyPattern == "transactionId:" {
				tx, err := wire.Svc.TransactionService.GetTransactionById(ctx, restOfString)
				if err != nil {
					continue
				}
				
			}

			// Perform action on the data (e.g., delete the key)
			err := redisClient.Del(ctx, key).Err()
			if err != nil {
				log.Println("Error deleting key", key, "from Redis:", err)
			} else {
				log.Println("Data with key", key, "is older than 50 seconds and has been processed.")
			}
		}
	}
}
