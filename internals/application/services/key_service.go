package services

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type KeyService struct{}

func (ks KeyService) GeneratePublicKey() (string, error) {
	uuid, err := uuid.NewUUID()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("PK_%s", uuid.String()), nil
}

func (ks KeyService) GenerateSecretKey() (string, error) {
	uuid, err := uuid.NewUUID()

	return uuid.String(), err
}

func (ks KeyService) IsPublicKey(key string) bool {
	return strings.HasPrefix(key, "PK_")
}

func NewKeyService() *KeyService {
	return &KeyService{}
}
