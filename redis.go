package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
)

// Flatten RedisWorkflowJob to map of strings
type RedisWorkflowToMap interface {
	structFlatten(value RedisWorkflowJob, tag string) (map[string]interface{}, error)
}

// Storage is an adapter between the application and storage implementation
type Storage interface {
	Put(value WorkflowJob)
}

// RedisStorage holds a redis connection pool
type RedisStorage struct {
	conn *redis.Client
	ctx  context.Context
}

// InitStorage initializes and returns connection to storage
func InitStorage(cfg Config) *RedisStorage {
	redisOptions, err := redis.ParseURL(cfg.RedisDSN)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	rdb := &RedisStorage{
		conn: redis.NewClient(redisOptions),
		ctx:  context.Background(),
	}

	return rdb
}

func (rdb *RedisStorage) Put(job WorkflowJob) {
	ttl, err := time.ParseDuration("5m")
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	ID := fmt.Sprintf("%s:%s", job.Repository, job.Name)

	runKey := fmt.Sprintf("runs:%s:%d", ID, job.RunID)
	redisJob, err := structFlatten(intoRedisWorkFlowJob(&job))
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	rdb.conn.HSet(rdb.ctx, runKey, redisJob)
	rdb.conn.Expire(rdb.ctx, runKey, ttl)

	statKey := fmt.Sprintf("stats:%s", ID)
	rdb.conn.SAdd(rdb.ctx, statKey)
	rdb.conn.Expire(rdb.ctx, statKey, ttl)

	if job.Conclusion == "success" {
		duration := job.CompletedAt.Sub(job.StartedAt).Seconds()
		rdb.conn.LPush(rdb.ctx, statKey, duration)
		rdb.conn.LTrim(rdb.ctx, statKey, 0, 99)
	}
}

func structFlatten(in interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("structFlatten only accepts struct or struct pointer")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if value := field.Tag.Get("redis"); value != "" {
			out[value] = v.Field(i).Interface()
		}
	}
	return out, nil
}
