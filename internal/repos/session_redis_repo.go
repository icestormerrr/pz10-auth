package repos

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionRedisRepoConfig struct {
	RefreshTTL    time.Duration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

type SessionRedisRepo struct {
	db     *redis.Client
	config SessionRedisRepoConfig
} // key = userID, value = refreshToken

func NewSessionRedisRepo(config SessionRedisRepoConfig) *SessionRedisRepo {
	db := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})
	return &SessionRedisRepo{db: db}
}

// TODO: Хранить «отозванные» refresh в in-memory blacklist (map) с exp.
func (repo *SessionRedisRepo) SetRefreshToken(userID int64, refreshTokenToSet string) error {
	return repo.db.Set(context.Background(), "auth/"+strconv.FormatInt(userID, 10), refreshTokenToSet, repo.config.RefreshTTL).Err()
}

func (repo *SessionRedisRepo) GetRefreshToken(userID int64) (string, error) {
	return repo.db.Get(context.Background(), "auth/"+strconv.FormatInt(userID, 10)).Result()
}
