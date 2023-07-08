package helpers

import (
	"robinhood-assignment/internal/core/ports"

	"golang.org/x/crypto/bcrypt"
)

type myBcrypt struct {
}

func NewMyBcrypt() ports.MyBcrypt {
	return &myBcrypt{}
}

func (b myBcrypt) GenerateFromPassword(password string, cost int) (*string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return nil, err
	}
	passHash := string(hashed)
	return &passHash, nil
}

func (b myBcrypt) CompareHashAndPassword(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}
