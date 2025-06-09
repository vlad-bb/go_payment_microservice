package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"time"
)

func (s *Service) CreateAuthToken(userID string) (string, error) {
	utcNow := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(utcNow),
		ExpiresAt: jwt.NewNumericDate(utcNow.Add(s.tokenTTLMinutes)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.cfg.AuthSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *Service) VerifyAuthToken(token string) (string, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.AuthSecret), nil
	})
	if err != nil {
		return "", err
	}
	if !jwtToken.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("failed to cast claims")
	}
	return claims.GetSubject()
}
