package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

const TABLE_USERS = `create table if not exists users
(
  id           serial,
  created_date timestamp default now(),
  chat_id      bigint unique     not null,
  access_token text unique     not null
);`

type Postgres struct {
	db *sql.DB
}

type User struct {
	id          int
	createdDate time.Time
	chatId      int
	accessToken string
}

func NewPostgres(url string) *Postgres {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(TABLE_USERS)
	if err != nil {
		log.Fatal(err)
	}

	return &Postgres{
		db,
	}
}

func (postgres *Postgres) Add(chatId int64, accessToken string) error {
	result, err := postgres.db.Exec("insert into users (chat_id, access_token) values ($1, $2)", chatId, accessToken)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("User %d created successfully (%d row affected)\n", chatId, rowsAffected)

	return nil
}

func (postgres *Postgres) Get(chatId int64) (string, error) {
	stmt, err := postgres.db.Prepare("select u.access_token from users u where u.chat_id=$1")
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	row := stmt.QueryRow(chatId)

	var accessToken string
	err = row.Scan(&accessToken)
	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (postgres *Postgres) Delete(chatId int64) error {
	result, err := postgres.db.Exec("delete from users where chat_id=$1", chatId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("User %d deleted successfully (%d row affected)\n", chatId, rowsAffected)

	return nil
}
