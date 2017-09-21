package storage

import (
	"sync"
	"fmt"
	"errors"
)

func NewTokenStore() *TokenStore {
	return &TokenStore{
		store: make(map[int64]string),
	}
}

type TokenStore struct {
	sync.Mutex
	store map[int64]string
}

func (s *TokenStore) Get(chatId int64) (string, error) {
	var token string
	s.Lock()
	token = s.store[chatId]
	s.Unlock()

	if token == "" {
		return "", errors.New(fmt.Sprintf("Token not found for chatId=%s\n", chatId))
	}
	return token, nil
}

func (s *TokenStore) Add(chatId int64, accessToken string) {
	s.Lock()
	s.store[chatId] = accessToken
	s.Unlock()
}

func (s *TokenStore) Delete(key int64) {
	s.Lock()
	delete(s.store, key)
	s.Unlock()
}
