package utils
import (
	"gopkg.in/redis.v3"
	"log"
)

// Everything that has to do with redis goes in here

func NewRedisClient() (*redis.Client) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0, // use default DB
	})
	return client
}

// writes result into redis DB to cache it.
// key: url, value: minified document
func writeToDB(client *redis.Client, newDBEntry map[string]string) (error) {
	if newDBEntry[""] == "cached" {
		return nil
	}

	for url, minifiedHTML := range newDBEntry {
		err := client.Set(url, minifiedHTML, 0).Err()
		if err != nil {
			log.Println("Failed to set value in redis db.")
		}
	}
	return nil
}

// queries one KEY of db, returns appropiate value if present
// or the empty string if not.
func queryDBKey(client *redis.Client, key string) (string) {
	res, err := client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return ""
		}
		log.Println("database error", err.Error())
		return ""
	}
	return res
}