package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	AuthHeader := headers.Get("Authorization")
	if AuthHeader == "" {
		return "", fmt.Errorf("no authorization header was found")
	}

	HeaderSplit := strings.Split(AuthHeader, " ")

	if len(HeaderSplit) > 2 || HeaderSplit[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return HeaderSplit[1], nil
}
