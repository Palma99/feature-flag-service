package services

import "golang.org/x/crypto/bcrypt"

type PasswordService struct{}

func (ps PasswordService) CompareHashAndPassword(rawPassword string, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword)) != nil
}

func (ps PasswordService) HashPassword(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func NewPasswordService() PasswordService {
	return PasswordService{}
}
