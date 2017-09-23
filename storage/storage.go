package storage

type AccessTokenStorage interface {
	Get(chatId int64) (string, error)
	Add(chatId int64, accessToken string) error
	Delete(key int64) error
}
