package choiceCache

import (
	"strconv"
	"time"

	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/go-redis/redis"
)

type choiceCache struct {
	logger *logging.Logger
	client *redis.Client
}

func NewChoiceCache(client *redis.Client, logger *logging.Logger) *choiceCache {
	return &choiceCache{logger: logger, client: client}
}

func (c *choiceCache) Set(voteTitle string, choiceTitle string, count int, expireAt time.Duration) error {
	c.logger.Infof("try to save %v : %v", voteTitle, choiceTitle)
	err := c.client.HSet(voteTitle, choiceTitle, count).Err()
	if err != nil {
		c.logger.Error(err)
		return err
	}
	err = c.client.Expire(voteTitle, expireAt).Err()
	if err != nil {
		c.logger.Error(err)
		return err
	}
	return nil
}

func (c *choiceCache) Get(voteTitle string, choiceTitle string) (int, error) {
	value, err := c.client.HGet(voteTitle, choiceTitle).Result()
	if err != nil {
		c.logger.Info(err)
		return -1, err
	}
	count, err := strconv.Atoi(value)
	if err != nil {
		c.logger.Info(err)
		return -1, err
	}
	return count, nil
}

func (c *choiceCache) GetAll(voteTitle string) (map[string]string, error) {
	m, err := c.client.HGetAll(voteTitle).Result()
	if err != nil {
		c.logger.Info(err)
		return nil, err
	}
	return m, nil
}
