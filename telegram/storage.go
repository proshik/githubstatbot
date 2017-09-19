package telegram

import (
	"sync"
	"errors"
	"fmt"
)

func NewStorage() Storage {
	return &Store{
		store: make(map[int64]string),
	}
}

type Storage interface {
	Get(chatId int64) (string, error)
	Add(chatId int64, accessToken string)
}

type Store struct {
	sync.Mutex
	store map[int64]string
}

func (s *Store) Get(chatId int64) (string, error) {
	var token string
	s.Lock()
	token = s.store[chatId]
	s.Unlock()

	if token == "" {
		return "", errors.New(fmt.Sprintf("Token not found for chatId=%s\n", chatId))
	}
	return token, nil
}

func (s *Store) Add(chatId int64, accessToken string) {
	s.Lock()
	s.store[chatId] = accessToken
	s.Unlock()
}
