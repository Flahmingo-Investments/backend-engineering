package tests

import (
	jwtService "grpc-example/src/commons/services/jwt"
	commons "grpc-example/src/commons/utils"
	"grpc-example/src/data/models"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	// Data
	user := &models.User{}
	user.ID = "3"

	// Service
	service := jwtService.JWTAuth{SECRET: "buu"}

	// Creates a token
	commons.TestingTime("00:10")
	token := service.GetJWT(user, jwt.SigningMethodHS256, 1)
	assert.NotEmpty(t, token)

	// Success Flow
	user_id := service.ParseJWT(token, jwt.SigningMethodHS256)
	assert.Equal(t, user_id, user.ID)

	// Validate signing method
	user_id = service.ParseJWT(token, jwt.SigningMethodHS512)
	assert.Empty(t, user_id)

	// Validate wrong secret
	service.SECRET = "no!"
	user_id = service.ParseJWT(token, jwt.SigningMethodHS256)
	assert.Empty(t, user_id)

	// Validate expired time
	commons.TestingTime("00:15")
	service.SECRET = "buu"
	user_id = service.ParseJWT(token, jwt.SigningMethodHS256)
	assert.Empty(t, user_id)
}
