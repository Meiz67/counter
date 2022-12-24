package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type Event struct {
	Id         string `json:"id"`
	RemoteAddr string `json:"remoteAddr"`
	Uri        string `json:"uri"`
	EventType  int    `json:"type"`
	Function   string `json:"function"`
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	err := client.Ping().Err()
	if err != nil {
		fmt.Println("Redis connection error!")
	}

	for {
		data, _ := client.BRPop(time.Second, "COUNTER").Result()
		var event Event
		if data != nil {
			//fmt.Println(data[1])
			err := json.Unmarshal([]byte(data[1]), &event)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(event.Id)
			}
		}
	}

}
