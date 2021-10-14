package tests

import (
	"context"
	"grpc-example/src/apps/auth/services"
	"grpc-example/src/apps/auth/tests/mocks"
	jwtAuth "grpc-example/src/commons/services/jwt"
	commons "grpc-example/src/commons/utils"
	"grpc-example/src/data/models"
	"grpc-example/src/proto"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestGetProfileWrongAccessToken(t *testing.T) {
	user := &models.User{}
	user.ID = "11"

	jwtService := jwtAuth.JWTAuth{SECRET: "buu"}
	authService := &services.AuthService{
		UserRepository: &mocks.UserRepositoryMock{
			GetResponse: user,
		},
		JWT: &jwtService,
	}
	commons.TestingTime("00:10")

	// Using the refresh token as access token
	response, _ := authService.GetProfile(context.Background(), &proto.ProfileRequest{
		AccessToken: jwtService.GetJWT(user, jwt.SigningMethodHS512, time.Minute*5),
	})
	assert.Equal(t, response.Status, proto.Status_UNAUTHORIZED)

	// Expired token
	token := jwtService.GetJWT(user, jwt.SigningMethodHS256, time.Minute*5)
	commons.TestingTime("00:20")
	response, _ = authService.GetProfile(context.Background(), &proto.ProfileRequest{
		AccessToken: token,
	})
	assert.Equal(t, response.Status, proto.Status_UNAUTHORIZED)

	// Using another secret key >(
	commons.TestingTime("00:10")
	jwtService.SECRET = "no!"
	response, _ = authService.GetProfile(context.Background(), &proto.ProfileRequest{
		AccessToken: token,
	})
	assert.Equal(t, response.Status, proto.Status_UNAUTHORIZED)
}

func TestGetProfileSuccessFlow(t *testing.T) {
	user := &models.User{}
	user.ID = "11"

	jwtService := jwtAuth.JWTAuth{SECRET: "buu"}
	authService := &services.AuthService{
		UserRepository: &mocks.UserRepositoryMock{
			GetResponse: user,
		},
		JWT: &jwtService,
	}
	commons.TestingTime("00:10")
	response, _ := authService.GetProfile(context.Background(), &proto.ProfileRequest{
		AccessToken: jwtService.GetJWT(user, jwt.SigningMethodHS256, time.Minute*5),
	})
	assert.Equal(t, response.Status, proto.Status_SUCCESS)
	assert.Equal(t, response.Id, user.ID)
}
