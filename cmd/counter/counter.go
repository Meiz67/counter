package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type Event struct {
	Id         string  `json:"id"`
	RemoteAddr string  `json:"remoteAddr"`
	Uri        string  `json:"uri"`
	EventType  int     `json:"type"`
	Function   string  `json:"function"`
	FuncTime   float32 `json:"funcTime"`
}

type Counter struct {
	Functions  map[string]FuncData
	ScriptTime float32
}

type FuncData struct {
	time  float32
	count int
}

func main() {
	Summary := map[string]Counter{}

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
		data, _ := client.BRPop(time.Second*10, "COUNTER").Result()
		var event Event
		if data != nil {
			err := json.Unmarshal([]byte(data[1]), &event)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				IncCounter(Summary, event.Id, event.Function)
			}
		}
		fmt.Println("Summary len: ", len(Summary))
		for key, value := range Summary {
			fmt.Println(key, value)
		}
	}

}

func IncCounter(summary map[string]Counter, id string, function string) {
	if _, ok := summary[id]; ok == true {
		if _, ok := summary[id].Functions[function]; ok == true {
			funcData := summary[id].Functions[function]
			funcData.count += 1
			summary[id].Functions[function] = funcData
		} else {
			summary[id].Functions[function] = FuncData{
				time:  0,
				count: 1,
			}
		}
	} else {
		summary[id] = Counter{
			Functions:  map[string]FuncData{},
			ScriptTime: 0,
		}
	}

}
