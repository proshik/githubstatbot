package storage

import (
	"sync"
	"fmt"
	"errors"
)

func NewStateStore() *StateStore {
	return &StateStore{
		store: make(map[string]int64),
	}
}

type StateStore struct {
	sync.Mutex
	store map[string]int64
}

func (s *StateStore) Get(state string) (int64, error) {
	var chatId int64
	s.Lock()
	chatId = s.store[state]
	s.Unlock()

	if chatId == 0 {
		return 0, errors.New(fmt.Sprintf("ChatId not found by state=%s\n", state))
	}
	return chatId, nil
}

func (s *StateStore) Add(state string, chatId int64) {
	s.Lock()
	s.store[state] = chatId
	s.Unlock()
}