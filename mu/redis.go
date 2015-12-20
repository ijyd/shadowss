package main

import (
	"gopkg.in/redis.v3"
	"github.com/orvice/shadowsocks-go/mu/log"
)

func InitRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil{
		log.Log.Error("init redis failed: ",err)
		return err
	}
	log.Log.Info(pong)
	return nil
}
