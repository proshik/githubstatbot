package storage

import (
	"encoding/binary"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

var tokenBucket = "token"

type Store struct {
	path string
}

func New(path string) *Store {
	db, err := bolt.Open(path, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(tokenBucket))
		return err
	})

	return &Store{
		path,
	}
}

func (s *Store) Add(chatId int64, accessToken string) error {
	db, err := open(s)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tokenBucket))
		err := b.Put(itob(chatId), []byte(accessToken))
		return err
	})
	if err != nil {
		log.Printf("error on add token, error: %s\n", err)
		return err
	}
	return nil
}

func (s *Store) Get(chatId int64) (string, error) {
	db, err := open(s)
	if err != nil {
		log.Printf("Error on get element for chatId=%s, %v\n", chatId, err)
		return "", err
	}
	defer db.Close()

	var token []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tokenBucket))
		token = b.Get(itob(chatId))
		return nil
	})

	return string(token), nil
}

func (s *Store) Delete(chatId int64) error {
	db, err := open(s)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tokenBucket))
		err := b.Delete(itob(chatId))
		return err
	})
	if err != nil {
		log.Printf("error on delete token, error: %s\n", err)
		return err
	}
	return nil
}

func open(s *Store) (*bolt.DB, error) {
	db, err := bolt.Open(s.path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Printf("Error on open database. %v\n", err)
		return nil, err
	}

	return db, nil
}

func itob(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
