package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/hovanhoa/go-url-shortener/config"
	"github.com/hovanhoa/go-url-shortener/internal/entities"
	"github.com/hovanhoa/go-url-shortener/internal/storage"
	"github.com/redis/go-redis/v9"
)

type (
	URLService interface {
		AddNewURL(url *entities.URL) (*entities.URL, error)
		FindOneByLongURL(longURL string) (*entities.URL, error)
		FindOneByShortURL(shortURL string) (*entities.URL, error)
		FindOneByID(id int64) (*entities.URL, error)
	}

	urlService struct {
		storage *storage.Repository
		redis   *redis.Client
	}
)

func (u *urlService) AddNewURL(url *entities.URL) (*entities.URL, error) {
	cfg := config.GetConfig()
	res, err := u.storage.URL.AddNewURL(url)
	if err != nil {
		return nil, err
	}

	err = u.redis.Set(context.Background(), res.SortURL, res.LongURL, cfg.Redis.ExpirationTime).Err()
	if err != nil {
		fmt.Println("Error on set key in redis, ", err)
	}

	return res, nil
}

func (u *urlService) FindOneByLongURL(longURL string) (*entities.URL, error) {
	url := entities.URL{
		LongURL: longURL,
	}

	res, err := u.storage.URL.FindOneURL(&url)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *urlService) FindOneByShortURL(shortURL string) (*entities.URL, error) {
	cfg := config.GetConfig()
	longURL, err := u.redis.Get(context.Background(), shortURL).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		fmt.Println("Error on get key in redis, ", err)
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

	res, err := u.storage.URL.FindOneURL(&url)
	if err != nil {
		return nil, err
	}

	err = u.redis.Set(context.Background(), res.SortURL, res.LongURL, cfg.Redis.ExpirationTime).Err()
	if err != nil {
		fmt.Println("Error on set key in redis, ", err)
	}

	return res, nil
}

func (u *urlService) FindOneByID(id int64) (*entities.URL, error) {
	url := entities.URL{
		ID: id,
	}

	res, err := u.storage.URL.FindOneURL(&url)
	if err != nil {
		return nil, err
	}

	return res, nil
}
