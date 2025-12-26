package dao

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

type Redis struct {
	redis *redis.Client
	ctx   context.Context
}

// 返回一个可用Redis

func NewRedis(redis *redis.Client, ctx context.Context) *Redis {
	return &Redis{redis: redis, ctx: ctx}
}

// 添加缓存过期机制

func (redis Redis) IsExpired(data Data) (bool, error) {
	ttl, err := redis.redis.TTL(context.Background(), data.Name).Result()
	if err != nil {
		return false, err
	}
	return ttl <= 0, nil
}

// 建立单个任务数据缓存

func (redis Redis) SetRedis(data Data, existTime time.Duration) error {
	result := redis.redis.HSet(redis.ctx, data.Name, data.TaskName, data.Value)
	err := redis.redis.Expire(context.Background(), data.Name, existTime).Err() //设置一个缓存时间
	if err != nil {
		return err
	}
	return result.Err()
}

// 得到任务数据

func (redis Redis) GetRedis(data Data) (interface{}, error) {
	result := redis.redis.HGet(redis.ctx, data.Name, data.TaskName)
	return result.Result()
}

// 建立多个任务数据缓存

func (redis Redis) GetRedisList(data Data) (interface{}, error) {
	result := redis.redis.HGetAll(redis.ctx, data.Name)
	return result.Result()
}

// 删除缓存

func (redis Redis) DeleteRedis(data Data) (bool, error) {
	result := redis.redis.HDel(redis.ctx, data.Name, data.TaskName)
	if result.Err() != nil {
		return false, result.Err()
	}
	if result.Val() == 0 {
		return false, nil
	}
	return true, result.Err()
}

// 分布式锁

func (redis Redis) Lock(data Data) (bool, error) {
	locKey := "lock:" + data.Name
	lockValue := strconv.FormatInt(time.Now().UnixNano(), 10)
	result, err := redis.redis.SetNX(redis.ctx, locKey, lockValue, 60*time.Second).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

//解锁

func (redis Redis) Unlock(data Data) error {
	locKey := "lock:" + data.Name
	result := redis.redis.Del(redis.ctx, locKey)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

//设置一个缓存储存内容

func (redis Redis) UpdateCache(data Data) error {
	key := data.Name + ":" + data.TaskName
	result := redis.redis.Set(redis.ctx, key, data.Value, 1*time.Hour)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

//设置一个备份缓存

func (redis Redis) AlwaysCache(data Data) error {
	key := "always" + data.Name + ":" + data.TaskName
	result := redis.redis.Set(redis.ctx, key, data.Value, 0)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

//从缓存提取数据

func (redis Redis) GetData(data string) (Data, error) {
	var reData Data
	result := redis.redis.Get(redis.ctx, data)
	if result.Err() != nil {
		return reData, result.Err()
	}
	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return reData, fmt.Errorf("错误的缓存格式")
	}
	reData.Name = parts[0]
	reData.TaskName = parts[1]
	reData.Value = result.Val()
	return reData, nil
}

//获得缓存内容

func (redis Redis) GetAll(data Data, cursor uint64) (uint64, []string, error) {
	match := "always" + data.Name + ":" + "*"
	result := redis.redis.Scan(redis.ctx, cursor, match, 10)
	if result.Err() != nil {
		return 0, nil, result.Err()
	}
	key, nextCursor := result.Val()
	return nextCursor, key, nil
}
