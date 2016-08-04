package main

import (
	"encoding/json"
	"time"

	muconfig "shadowsocks/shadowsocks-go/mu/config"
	"shadowsocks/shadowsocks-go/mu/user"
	redis "gopkg.in/redis.v3"
)

const (
	DefaultExpireTime          = 0
	DefaultOnlineKeyExpireTime = time.Minute * 5
)

var Redis = new(RedisClient)

type RedisClient struct {
	client *redis.Client
}

func (r *RedisClient) SetClient(client *redis.Client) {
	r.client = client
}

func (r *RedisClient) GetUserInfo(u user.User) (user.UserInfo, error) {
	var user user.UserInfo
	val, err := r.client.Get(genUserInfoKey(u.GetUserInfo())).Result()
	if err != nil {
		return user, err
	}
	err = json.Unmarshal([]byte(val), &user)
	return user, err
}

func (r *RedisClient) StoreUser(user user.UserInfo) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = r.client.Set(genUserInfoKey(user), data, DefaultExpireTime).Err()
	return err
}

func (r *RedisClient) Exists(u user.User) (bool, error) {
	return r.client.Exists(genUserInfoKey(u.GetUserInfo())).Result()
}

func (r *RedisClient) Del(u user.User) error {
	return r.client.Del(genUserInfoKey(u.GetUserInfo())).Err()
}

func (r *RedisClient) ClearAll() error {
	return r.client.FlushAll().Err()
}

// traffic
func (r *RedisClient) IncrSize(u user.User, size int) error {
	key := genUserFlowKey(u.GetUserInfo())
	incrSize := int(float32(size))
	isExits, err := r.client.Exists(key).Result()
	if err != nil {
		return err
	}
	if !isExits {
		return r.client.Set(key, incrSize, DefaultExpireTime).Err()
	}
	return r.client.IncrBy(key, int64(incrSize)).Err()
}

func (r *RedisClient) GetSize(u user.User) (int64, error) {
	key := genUserFlowKey(u.GetUserInfo())
	isExits, err := r.client.Exists(key).Result()
	if err != nil {
		return 0, err
	}
	if !isExits {
		return 0, nil
	}
	return r.client.Get(key).Int64()
}

func (r *RedisClient) SetSize(u user.User, size int) error {
	key := genUserFlowKey(u.GetUserInfo())
	return r.client.Set(key, size, DefaultExpireTime).Err()
}

func (r *RedisClient) ClearSize() error {
	return nil
}

func (r *RedisClient) MarkUserOnline(u user.User) error {
	key := genUserOnlineKey(u.GetUserInfo())
	return r.client.Set(key, "1", DefaultOnlineKeyExpireTime).Err()
}

func (r *RedisClient) IsUserOnline(u user.User) bool {
	key := genUserOnlineKey(u.GetUserInfo())
	isExits, err := r.client.Exists(key).Result()
	if err != nil {
		return false
	}
	return isExits
}

func (r *RedisClient) GetOnlineUsersCount(users []user.User) int {
	count := 0
	for _, v := range users {
		if r.IsUserOnline(v) {
			count++
		}
	}
	return count
}

func InitRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr:     muconfig.Conf.Redis.Host,
		Password: muconfig.Conf.Redis.Pass, // no password set
		DB:       muconfig.Conf.Redis.Db,   // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		return err
	}
	Log.Info(pong)
	Redis.SetClient(client)
	// set storage
	SetStorage(Redis)
	return nil
}
