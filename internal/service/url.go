package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/hovanhoa/go-url-shortener/config"
	"github.com/hovanhoa/go-url-shortener/internal/entities"
	"github.com/hovanhoa/go-url-shortener/internal/storage"
	"github.com/hovanhoa/go-url-shortener/pkg/otel"
	"github.com/redis/go-redis/v9"
)

type (
	URLService interface {
		AddNewURL(ctx context.Context, url *entities.URL) (*entities.URL, error)
		FindOneByLongURL(ctx context.Context, longURL string) (*entities.URL, error)
		FindOneByShortURL(ctx context.Context, shortURL string) (*entities.URL, error)
		FindOneByID(ctx context.Context, id int64) (*entities.URL, error)
	}

	urlService struct {
		storage *storage.Repository
		redis   *redis.Client
	}
)

func (u *urlService) AddNewURL(ctx context.Context, url *entities.URL) (*entities.URL, error) {
	cfg := config.GetConfig()
	res, err := u.storage.URL.AddNewURL(ctx, url)
	if err != nil {
		return nil, err
	}

	err = otel.RedisWithTracing(ctx, "set", res.SortURL, func() error {
		return u.redis.Set(ctx, res.SortURL, res.LongURL, cfg.Redis.ExpirationTime).Err()
	})
	if err != nil {
		slog.Warn("Failed to set key in Redis", "key", res.SortURL, "error", err)
	}

	return res, nil
}

func (u *urlService) FindOneByLongURL(ctx context.Context, longURL string) (*entities.URL, error) {
	url := entities.URL{
		LongURL: longURL,
	}

	res, err := u.storage.URL.FindOneURL(ctx, &url)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *urlService) FindOneByShortURL(ctx context.Context, shortURL string) (*entities.URL, error) {
	cfg := config.GetConfig()
	var longURL string
	err := otel.RedisWithTracing(ctx, "get", shortURL, func() error {
		var redisErr error
		longURL, redisErr = u.redis.Get(ctx, shortURL).Result()
		return redisErr
	})
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.Warn("Failed to get key from Redis", "key", shortURL, "error", err)
	}

	if longURL != "" {
		return &entities.URL{
			SortURL: shortURL,
			LongURL: longURL,
		}, nil
	}

	url := entities.URL{
		SortURL: shortURL,
	}

	res, err := u.storage.URL.FindOneURL(ctx, &url)
	if err != nil {
		return nil, err
	}

	err = otel.RedisWithTracing(ctx, "set", res.SortURL, func() error {
		return u.redis.Set(ctx, res.SortURL, res.LongURL, cfg.Redis.ExpirationTime).Err()
	})
	if err != nil {
		slog.Warn("Failed to set key in Redis", "key", res.SortURL, "error", err)
	}

	return res, nil
}

func (u *urlService) FindOneByID(ctx context.Context, id int64) (*entities.URL, error) {
	url := entities.URL{
		ID: id,
	}

	res, err := u.storage.URL.FindOneURL(ctx, &url)
	if err != nil {
		return nil, err
	}

	return res, nil
}
