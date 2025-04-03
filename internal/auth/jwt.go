package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	stringToUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}

	return stringToUUID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	AuthHeader := headers.Get("Authorization")
	if AuthHeader == "" {
		return "", fmt.Errorf("no authorization header was found")
	}

	HeaderSplit := strings.Split(AuthHeader, " ")

	if len(HeaderSplit) > 2 || HeaderSplit[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return HeaderSplit[1], nil
}
