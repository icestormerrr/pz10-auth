package repos

import (
	"context"
	"strconv"

	"github.com/icestormerrr/pz10-auth/internal/utils/config"
	"github.com/redis/go-redis/v9"
)

type SessionRedisRepo struct {
	db     *redis.Client
	config config.Config
} // key = userID, value = refreshToken

func NewSessionRedisRepo(config config.Config) *SessionRedisRepo {
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
