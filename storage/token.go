package storage

import (
	"sync"
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

func (s *TokenStore) Get(chatId int64) string {
	s.Lock()
	defer s.Unlock()
	return s.store[chatId]
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
