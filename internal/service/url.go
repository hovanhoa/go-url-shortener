package service

import (
	"github.com/hovanhoa/go-url-shortener/internal/entities"
	"github.com/hovanhoa/go-url-shortener/internal/storage"
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
	}
)

func (u *urlService) AddNewURL(url *entities.URL) (*entities.URL, error) {
	res, err := u.storage.URL.AddNewURL(url)
	if err != nil {
		return nil, err
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
	url := entities.URL{
		SortURL: shortURL,
	}

	res, err := u.storage.URL.FindOneURL(&url)
	if err != nil {
		return nil, err
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
