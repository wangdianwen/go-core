package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenJWT(expireTime int64, secretKeyStr string) (ss string) {
	if expireTime == 0 {
		expireTime = 10
	}
	secretKey := make([]byte, 0)
	if len(secretKeyStr) > 0 {
		secretKey = []byte(secretKeyStr)
	}
	claims := &jwt.StandardClaims{
		Audience:  "Brunton Applications",
		ExpiresAt: time.Now().Unix() + expireTime,
		Id:        "brunton",
		IssuedAt:  time.Now().Unix(),
		Issuer:    "brunton.co.nz",
		NotBefore: time.Now().Unix(),
		Subject:   "Brunton Inner Service",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return ""
	}
	return ss
}

func CheckJWT(secretKeyStr string, token string) (bool, error) {
	secretKey := make([]byte, 0)
	if len(secretKeyStr) > 0 {
		secretKey = []byte(secretKeyStr)
	}
	if len(token) == 0 {
		return false, errors.New("invalid parameter: token")
	}
	t, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if t != nil && t.Valid {
		return true, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false, errors.New("invalid jwt token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return false, errors.New("invalid jwt token, please check your system time")
		} else {
			return false, errors.New("invalid jwt token")
		}
	} else {
		return false, errors.New("invalid jwt token")
	}
}
