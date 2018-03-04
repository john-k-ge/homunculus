package cs

import (
	"errors"
	"log"
	"strconv"
	"time"

	"gopkg.in/redis.v3"
)

type RemoteCS struct {
	redisClient *redis.Client
	suffix      string
}

func NewRemoteCondtionSet(host, port, pass string, index int) (*RemoteCS, error) {
	temp := RemoteCS{
		redisClient: redis.NewClient(
			&redis.Options{
				Addr:     host + ":" + port,
				Password: pass,
				DB:       0,
			}),
		suffix: strconv.Itoa(index),
	}

	_, err := temp.redisClient.Ping().Result()
	if err != nil {
		log.Printf("New Redis client ping failed: %v", err)
		return nil, err
	}
	return &temp, nil
}

func (c *RemoteCS) SaveCondition(condition string, count int64) error {
	if len(condition) == 0 {
		return errors.New("must supply a key")
	}
	return c.put(condition+c.suffix, count)
}

func (c *RemoteCS) CheckCondition(condition string) (int64, error) {
	return c.get(condition + c.suffix)
}

func (c *RemoteCS) ConditionExists(condition string) (bool, error) {
	return c.exists(condition + c.suffix)
}

func (c *RemoteCS) DeleteCondition(condition string) error {
	return c.redisClient.Del(condition + c.suffix).Err()
}

func (c *RemoteCS) ClearConditions() error {
	return c.redisClient.FlushAll().Err()
}

func (c *RemoteCS) IncrementCondition(condition string) (int64, error) {
	return c.incr(condition + c.suffix)
}

func (c *RemoteCS) Size() (int64, error) {
	return c.redisClient.DbSize().Result()
}

// Private functions

func (c *RemoteCS) put(k string, v int64) error {
	if len(k) == 0 {
		log.Printf("Either key or value nil -- %v : %v", k, v)
		return errors.New("either key or value nil")
	}

	err := c.redisClient.Set(k, v, time.Hour).Err()
	if err != nil {
		log.Printf("Failed to set cache entry (%v:%v):  %v", k, v, err)
		return err
	}
	return nil
}

func (c *RemoteCS) get(k string) (int64, error) {
	if len(k) == 0 {
		log.Print("Key is empty")
		return 0, errors.New("key is empty")
	}
	v, err := c.redisClient.Get(k).Result()
	if err != nil {
		log.Printf("Failed to get value for key %v:  %v", k, err)
		return 0, err
	}
	log.Printf("For k=%v, found v=%v", k, v)
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("Failed to convert %v to int: %v", v, err)
	}
	return int64(i), nil
}

func (c *RemoteCS) exists(k string) (bool, error) {
	if len(k) == 0 {
		log.Print("Key is empty")
		return false, errors.New("key is empty")
	}
	return c.redisClient.Exists(k).Result()
}

func (c *RemoteCS) incr(k string) (int64, error) {
	if len(k) == 0 {
		log.Print("Key is empty")
		return 0, errors.New("key is empty")
	}

	res, err := c.redisClient.Incr(k).Result()
	log.Printf("result: %v", res)

	return res, err
}
