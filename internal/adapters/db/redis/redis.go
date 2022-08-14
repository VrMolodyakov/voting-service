package redis

import (
	"encoding/json"
	"time"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/go-redis/redis"
)

type choiceCache struct {
	logger *logging.Logger
	client *redis.Client
}

func NewRedisClient(client *redis.Client, logger *logging.Logger) *choiceCache {
	return &choiceCache{logger: logger, client: client}
}

func (c *choiceCache) Save(voteTitle string, choice entity.Choice, expireAt time.Duration) error {
	json, err := json.Marshal(choice)
	if err != nil {
		return err
	}
	return c.client.Set(voteTitle, json, expireAt).Err()
}

func (c *choiceCache) Get(voteTitle string) (entity.Choice, error) {
	str, err := c.client.Get(voteTitle).Result()
	if err != nil {
		return entity.Choice{}, err
	}
	var choice entity.Choice
	err = json.Unmarshal([]byte(str), &choice)
	return choice, err
}
