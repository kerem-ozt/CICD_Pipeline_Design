package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB() {
	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, Config.MongodbDatabase, options.Client().ApplyURI(Config.MongodbUri))
	if err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB!")
}

var redisDefaultClient *redis.Client
var redisDefaultOnce sync.Once

var redisCache *cache.Cache
var redisCacheOnce sync.Once

func GetRedisDefaultClient() *redis.Client {
	redisDefaultOnce.Do(func() {
		redisDefaultClient = redis.NewClient(&redis.Options{
			Addr: Config.RedisDefaultAddr,
		})
	})

	return redisDefaultClient
}

func GetRedisCache() *cache.Cache {
	redisCacheOnce.Do(func() {
		redisCache = cache.New(&cache.Options{
			Redis:      GetRedisDefaultClient(),
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		})
	})

	return redisCache
}

func CheckRedisConnection() {
	redisClient := GetRedisDefaultClient()
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	log.Println("Connected to Redis!")
}

func getTournamentCacheKey(tournamentId primitive.ObjectID) string {
	return "tournament:" + tournamentId.Hex()
}

func CacheOneTournament(tournamentId primitive.ObjectID, winners []Participant) error {
	if !Config.UseRedis {
		return nil
	}

	// tournamentCacheKey := getTournamentCacheKey(userId, tournament.ID)

	tournamentCacheKey := "tournament:" + tournamentId.Hex()

	// Marshal the tournament data to JSON before caching
	tournamentJSON, err := json.Marshal(winners)

	if err != nil {
		return err
	}

	// Create a cache item and set it in Redis
	item := &cache.Item{
		Ctx:   context.TODO(),
		Key:   tournamentCacheKey,
		Value: tournamentJSON,
		TTL:   time.Minute,
	}

	return GetRedisCache().Set(item)
}

func GetTournamentFromCache(tournamentId primitive.ObjectID) (string, error) {
	if !Config.UseRedis {
		return "", errors.New("no redis client, set USE_REDIS in .env")
	}

	tournamentCacheKey := getTournamentCacheKey(tournamentId)

	var result string
	err := GetRedisCache().Get(context.TODO(), tournamentCacheKey, &result)
	if err != nil {
		fmt.Println("err:", err)
		return "", err
	}

	return result, nil
}
