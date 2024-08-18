package auth

import (
	"errors"
	"net/http"
)

func GetEmail(headers http.Header) (string, error) {
	val := headers.Get("Email")
	if val == "" {
		return "", errors.New("no email found")
	}
	return val, nil
}

func GetPassword(headers http.Header) (string, error) {
	val := headers.Get("Password")
	if val == "" {
		return "", errors.New("no password found")
	}
	return val, nil
}
