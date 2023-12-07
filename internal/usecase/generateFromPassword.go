package usecase

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateFromPassword(password string) ([]byte, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashPassword, nil
}
