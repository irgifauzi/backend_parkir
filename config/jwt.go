package config

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
)

var JWTSecret = "ZeNvErSERynHrSZ"

func GenerateJWT(adminID string) (string, error) {
    claims := jwt.MapClaims{
        "admin_id": adminID,
        "exp":      time.Now().Add(time.Hour * 12).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(JWTSecret))
}