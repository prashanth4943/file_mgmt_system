package storage

import (
	"context"
	"file_mgmt_system/helper"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisClient(addr string, password string, db int) (*RedisClient, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Redis successfully")
	return &RedisClient{
		Client: client,
		Ctx:    ctx,
	}, nil
}

func (r *RedisClient) Close() {
	if err := r.Client.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}
}

func (r *RedisClient) SetFile(fileID string, fileData []byte) error {
	compressedData, err := helper.Compress(fileData)
	if err != nil {
		log.Printf("Error compressing file: %v", err)
		return err
	}

	err = r.Client.Set(r.Ctx, "file:"+fileID, compressedData, time.Hour).Err()
	if err != nil {
		log.Printf("Error storing file in Redis: %v", err)
		return err
	}

	log.Printf("File %s stored successfully in Redis", fileID)
	return nil
}

func (r *RedisClient) GetFile(fileID string) ([]byte, error) {
	compressedData, err := r.Client.Get(r.Ctx, "file:"+fileID).Bytes()
	if err != nil {
		log.Printf("Error retrieving file from Redis: %v", err)
		return nil, err
	}

	fileData, err := helper.Decompress(compressedData)
	if err != nil {
		log.Printf("Error decompressing file: %v", err)
		return nil, err
	}

	log.Printf("File %s retrieved successfully from Redis", fileID)
	return fileData, nil
}

func (r *RedisClient) DeleteFile(fileID string) error {
	err := r.Client.Del(r.Ctx, "file:"+fileID).Err()
	if err != nil {
		log.Printf("Error deleting file from Redis: %v", err)
		return err
	}

	log.Printf("File %s deleted successfully from Redis", fileID)
	return nil
}
