package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct{}

func (ps PasswordService) ArePasswordsEqual(rawPassword string, hashedPassword string) bool {
	result := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))

	return result == nil
}

func (ps PasswordService) HashPassword(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

func isValidSymbol(char rune) bool {
	validSymbols := []rune{'_', '-', '.'}

	for _, symbol := range validSymbols {
		if char == symbol {
			return true
		}
	}

	return false
}

func (ps PasswordService) CheckPasswordSecurity(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if len(password) > 32 {
		return errors.New("password must be at most 32 characters long")
	}

	symbols := 0
	uppercase := 0
	lowercase := 0
	numbers := 0

	for _, char := range password {
		switch true {
		case char >= 'a' && char <= 'z':
			lowercase++
		case char >= 'A' && char <= 'Z':
			uppercase++
		case isValidSymbol(char):
			symbols++
		case char >= '0' && char <= '9':
			numbers++
		default:
			return errors.New("password must contain only letters, numbers and symbols (_, -, .)")
		}

	}

	if symbols == 0 {
		return errors.New("password must contain at least one symbol")
	}

	if uppercase == 0 {
		return errors.New("password must contain at least one uppercase letter")
	}

	if lowercase == 0 {
		return errors.New("password must contain at least one lowercase letter")
	}

	if numbers == 0 {
		return errors.New("password must contain at least one number")
	}

	return nil
}
