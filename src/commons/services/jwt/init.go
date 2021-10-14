package jwt

import (
	"grpc-example/src/data/models"
	"os"
	"time"

	jwtLib "github.com/golang-jwt/jwt"
)

type JWTManager interface {
	GetJWT(*models.User, jwtLib.SigningMethod, time.Duration) string
	ParseJWT(string, jwtLib.SigningMethod) string
}

type JWTAuth struct {
	SECRET string
}

func Initialize() *JWTAuth {
	return &JWTAuth{
		SECRET: os.Getenv("JWT_SECRET"),
	}
}
