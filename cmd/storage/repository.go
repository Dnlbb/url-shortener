package storage

type Repository interface {
	Save(shortURL, originalURL string) error
	Find(shortURL string) (string, bool)
}
