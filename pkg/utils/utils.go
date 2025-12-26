// utils/jwt.go
package utils

import (
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
)

type MyClaims struct {
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

func MakeToken(name string) (string, error) {
	myClaims := &MyClaims{
		UserName: name,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "mianmianyixia",
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
	token, err := t.SignedString([]byte("vivo50"))
	if err != nil {
		return "", err
	}
	return token, nil
}

func RandomDuration(rt int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randTime := time.Duration(r.Intn(rt)) * time.Hour
	return randTime
}
