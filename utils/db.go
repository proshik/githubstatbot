package db

import (
	"github.com/boltdb/bolt"
	"strconv"
	"time"
)

var tokenBucket = "token"

type Store struct {
	path string
}

func New(path string) *Store {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(tokenBucket))
		return err
	})

	return &Store{
		path,
	}
}

func (s *Store) Add(chatId int64, accessToken string) {
	db, _ := open(s)
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tokenBucket))
		err := b.Put([]byte(strconv.Itoa(int(chatId))), []byte(accessToken))
		return err
	})
}

func (s *Store) Get(chatId int64) string {
	db, _ := open(s)
	defer db.Close()

	var token []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tokenBucket))
		token = b.Get([]byte(strconv.Itoa(int(chatId))))
		return nil
	})

	return string(token)
}

func (s *Store) Delete(chatId int64) {
	db, _ := open(s)
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tokenBucket))
		err := b.Delete([]byte(strconv.Itoa(int(chatId))))
		return err
	})
}

func open(s *Store) (*bolt.DB, error) {
	db, err := bolt.Open(s.path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return db, nil
}
