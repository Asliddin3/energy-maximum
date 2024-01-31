package hash

import "golang.org/x/crypto/bcrypt"

type Password interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, check string) error
}

type Hash struct{}

func NewHasher() *Hash {
	return &Hash{}
}

func (h *Hash) HashPassword(password string) (string, error) {
	pw := []byte(password)
	result, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (h *Hash) CheckPassword(password, check string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(check))
	if err != nil {
		return err
	}
	return nil
}
