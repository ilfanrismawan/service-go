package utils

import (
    "context"
    "time"
    "service/internal/database"
)

const refreshBlacklistPrefix = "refresh_blacklist:"

// RevokeRefreshToken stores the refresh token in blacklist with TTL until expiry
func RevokeRefreshToken(ctx context.Context, token string, expiresAt time.Time) error {
    if database.Redis == nil {
        return nil
    }
    ttl := time.Until(expiresAt)
    if ttl <= 0 {
        ttl = time.Hour // fallback minimal
    }
    key := refreshBlacklistPrefix + token
    return database.Redis.Set(ctx, key, "1", ttl).Err()
}

// IsRefreshTokenRevoked checks if the refresh token is blacklisted
func IsRefreshTokenRevoked(ctx context.Context, token string) (bool, error) {
    if database.Redis == nil {
        return false, nil
    }
    key := refreshBlacklistPrefix + token
    res, err := database.Redis.Exists(ctx, key).Result()
    if err != nil {
        return false, err
    }
    return res == 1, nil
}


