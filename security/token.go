package security

import (
	"fibergo/utils"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

var (
	SecretKey     = []byte(os.Getenv("SECRETJWTKEY"))
	SigningMethod = jwt.SigningMethodHS256.Name
)

func NewToken(userId string) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		Id:        userId,
		IssuedAt:  time.Now().Unix(),
		Issuer:    userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func validateSignedMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return SecretKey, nil
}

func ParseToken(tokenString string) (*jwt.StandardClaims, error) {
	claims := new(jwt.StandardClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, validateSignedMethod)
	if err != nil {
		return nil, err
	}
	var ok bool
	claims, ok = token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return nil, utils.ErrInvalidAuthToken
	}
	return claims, nil
}
