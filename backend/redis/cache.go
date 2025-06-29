package redis

import (
	"backend/forest"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

func GuardarModelo(ctx context.Context, rdb *redis.Client, clave string, bosque forest.Bosque) error {
	bytes, err := json.Marshal(bosque)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, clave, bytes, 5*time.Minute).Err() // expira en 5 min
}

func LeerModelo(ctx context.Context, rdb *redis.Client, clave string) (forest.Bosque, error) {
	val, err := rdb.Get(ctx, clave).Bytes()
	if err != nil {
		return forest.Bosque{}, err
	}
	return forest.FromJSON(val)
}