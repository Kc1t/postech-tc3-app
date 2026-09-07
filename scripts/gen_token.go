package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func main() {
	claims := jwt.MapClaims{
		"sub": "admin-postman",
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("change-me-in-production"))
	fmt.Println(signed)
}
