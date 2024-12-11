package services

import (
	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	UserID   int
	Nickname string
}

type JWTService struct {
	secret string
}

func NewJWTService(secret string) *JWTService {
	if secret == "" {
		panic("secret is empty")
	}

	return &JWTService{
		secret,
	}
}

func (js JWTService) GenerateToken(userID int, nickname string) (string, error) {
	payload := &Payload{
		UserID:   userID,
		Nickname: nickname,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   payload.UserID,
		"nickname": payload.Nickname,
	})

	tokenString, err := token.SignedString([]byte(js.secret))

	return tokenString, err
}

func (js JWTService) ValidateToken(token string) (*Payload, error) {

	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(js.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(jwt.MapClaims); ok && tokenClaims.Valid {
		userID := int(claims["userID"].(float64))
		p := &Payload{UserID: userID, Nickname: claims["nickname"].(string)}
		return p, nil
	}

	return nil, err
}
