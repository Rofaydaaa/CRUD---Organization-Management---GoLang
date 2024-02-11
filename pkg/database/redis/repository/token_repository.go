package repository_token

import (
    "context"
    "fmt"

    "github.com/go-redis/redis/v8"
)

type TokenRepository struct {
    RedisClient *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) *TokenRepository {
    return &TokenRepository{RedisClient: redisClient}
}

func (repo *TokenRepository) SaveRefreshToken(userID, refreshToken string) error {
    ctx := context.Background()
    key := fmt.Sprintf("refresh_token:%s", userID)
    err := repo.RedisClient.Set(ctx, key, refreshToken, 0).Err()
    if err != nil {
        return err
    }
    return nil
}

func (repo *TokenRepository) RevokeRefreshToken(refreshToken string) error {
    ctx := context.Background()
    // Get the user ID associated with the refresh token from Redis
    userID, err := repo.RedisClient.Get(ctx, refreshToken).Result()
    if err != nil {
        return err
    }
    // Construct the key for the refresh token
    key := fmt.Sprintf("refresh_token:%s", userID)
    // Delete the refresh token from Redis
    err = repo.RedisClient.Del(ctx, key).Err()
    if err != nil {
        return err
    }
    return nil
}

func (repo *TokenRepository) RevokeRefreshTokenWithId(userID string) error {
    ctx := context.Background()
    key := fmt.Sprintf("refresh_token:%s", userID)
    err := repo.RedisClient.Del(ctx, key).Err()
    if err != nil {
        return err
    }
    return nil
}
