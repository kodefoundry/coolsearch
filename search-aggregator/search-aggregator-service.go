/*
	This service aggregates the results from individual search engines.
	The individual search results are posted on to a Kafka message broker on 'searchengine-result' topic
	Each message has the below structure

	{
		'searchengine': 'Name of search engine',
		'uuid': 'Unique id for the search',
		'keyword': 'Search key word'
		'results': [
			{
				'title': 'Title of the result',
				'link': 'Link to the original page',
				'snippet': 'Snippet text'
			}
		]
	}

	The aggregator does the following while consuming each message read from the queue
		1. Extract the uuid from message and search in its cache (Redis)
			2. If the uuid is not found create a new item in cache with uuid as key
				Format of Redis cache item
				uuid:{
					keyword: keyword
					google: $results (from above),
					ddg: $results (from above),
					wiki: $results (from above)
				}
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// SearchEngineResult stores the aggregated result from Google, Duck Duck Go and Wiki for a search keyword
type searchEngineResult struct {
	Searchengine string `json:"searchengine"`
	UUID         string `json:"uuid"`
	Keyword      string `json:"keyword"`
	Result       string `json:"results"`
}

type redisData struct {
	NoOfResult int
	Results    string
}

func (rd *redisData) addResult(result string) {
	rd.NoOfResult++
	rd.Results += "," + result
}

func main() {

	// Prepare Kafka Consumer
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "search-aggregator",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	// Prepare a redis client
	redisClient := redisClient()

	// Subscribe to 'searchengine-result' Kafka topic
	kafkaConsumer.SubscribeTopics([]string{"searchengine-result"}, nil)

	// Run an infinite loop and listen for incoming messages
	// When new message comes process the message by caching in redis
	// until all three search engine results are recieved for a serach id
	for {
		msg, err := kafkaConsumer.ReadMessage(-1)
		if err == nil {

			// uuid - int, searchResult - Json
			uuid, keyword, searchResult, searchengine := processMessage(string(msg.Value))
			fmt.Printf("Message received for %s, engine - %s\n", string(uuid), searchengine)
			// fetch the cached key (uuid) from redis, if not present insert new
			// returns a cachedResult - map, nRes - int, err
			nRes, cachedResult, err := getCachedResult(uuid, redisClient)
			if err != nil {
				fmt.Printf("Search result not found in cache.\n")
			}
			// returns err if uuid not found in cache ie the first result for a new search, add to cache
			switch {
			case nRes == 0:
				data := "1{\"searchkey\":\"" + keyword + "\"," + searchResult + ","
				redisClient.Set(uuid, data, 0)
			case nRes == 1:
				cachedResult = "2" + cachedResult[1:]
				redisClient.Set(uuid, cachedResult+searchResult+",", 0)
				break
			case nRes == 2:
				// Prepare the search result js
				data := cachedResult[1:] + searchResult + "}"
				redisClient.Set(uuid, data, 0)
				print("Posting data to 'search-result'\n")
				data = "{\"uuid\":\"" + uuid + "\", \"searchresult\":\"" + data + "\"}"
				postSearchResult("search-result", data)
				fmt.Print("Notifying search result availability to 'search-available'\n")
				postSearchResult("search-available", uuid)
				fmt.Print("Finished processing message-", uuid, "\n")
			}

		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}

	}
}

func redisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		print("Failed to connect to redis store, error - ", err, "\n")
	} else {
		print("Connected to redis store, ping result - ", pong, "\n")
	}
	return client
}

func postSearchResult(topic string, data string) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(data),
	}, nil)
}

/*
	{
		'searchengine': 'Name of search engine',
		'uuid': 'Unique id for the search',
		'keyword': 'Search key word'
		'results': [
			{
				'title': 'Title of the result',
				'link': 'Link to the original page',
				'snippet': 'Snippet text'
			}
		]
	}
*/
// This method process messages of above structure from the kafka topic 'search'
func processMessage(msg string) (uuid, keyWord, searchResult, searchengine string) {
	result := searchEngineResult{}
	json.Unmarshal([]byte(msg), &result)
	uuid = result.UUID
	keyWord = result.Keyword
	searchengine = result.Searchengine
	searchResult = "\"" + searchengine + "\":" + result.Result
	return uuid, keyWord, searchResult, searchengine
}

/*
uuid:{
	keyword: keyword,
	google: $results (from above),
	ddg: $results (from above),
	wiki: $results (from above)
}
*/
// This method manages data out from redis cache in the above format
func getCachedResult(uuid string, client *redis.Client) (int, string, error) {
	cachedResult, err := client.Get(uuid).Result()
	rNo := 0
	if err != nil || len(cachedResult) == 0 {
		print("No cached data found for id - ", uuid, "\n")
	} else {
		rNo, err = strconv.Atoi(cachedResult[:1])
	}
	return rNo, cachedResult, err
}
