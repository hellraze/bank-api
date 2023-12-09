package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateFromPassword(password []byte) ([]byte, error) {
	hashPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashPassword, nil
}

func CompareHashAndPassword(password []byte, hashPassword []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(password, hashPassword)
	if err != nil {
		return false, err
	}
	return true, err
}
