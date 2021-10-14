package tests

import (
	"context"
	"grpc-example/src/apps/auth/services"
	"grpc-example/src/apps/auth/tests/mocks"
	"grpc-example/src/data/models"
	"grpc-example/src/proto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignupDataValidation(t *testing.T) {
	authService := services.AuthService{}
	response, _ := authService.SignupWithPhoneNumber(context.Background(), &proto.SignupRequest{
		Name:  "",
		Phone: "0123456789",
	})
	assert.Equal(t, response.Status, proto.Status_ERROR)
	assert.Equal(t, *response.Detail, "Empty name")

	response, _ = authService.SignupWithPhoneNumber(context.Background(), &proto.SignupRequest{
		Name:  "Sergio",
		Phone: "",
	})
	assert.Equal(t, response.Status, proto.Status_ERROR)
	assert.Equal(t, *response.Detail, "Empty phone")

	response, _ = authService.SignupWithPhoneNumber(context.Background(), &proto.SignupRequest{
		Name:  "Sergio",
		Phone: "1234",
	})
	assert.Equal(t, response.Status, proto.Status_ERROR)
	assert.Equal(t, *response.Detail, "Phone number structure: +1012345789")
}
func TestSignupErrorInDatabase(t *testing.T) {
	authService := &services.AuthService{
		UserRepository: &mocks.UserRepositoryMock{
			CreateResponse: false,
		},
	}
	response, _ := authService.SignupWithPhoneNumber(context.Background(), &proto.SignupRequest{
		Name:  "Sergio",
		Phone: "+11234123412",
	})
	assert.Equal(t, response.Status, proto.Status_ERROR)
	assert.Equal(t, *response.Detail, "Error creating user on database")
}

func TestSignupUserExists(t *testing.T) {
	authService := &services.AuthService{
		UserRepository: &mocks.UserRepositoryMock{
			CreateResponse: true,
			GetResponse:    &models.User{},
		},
	}
	response, _ := authService.SignupWithPhoneNumber(context.Background(), &proto.SignupRequest{
		Name:  "Sergio",
		Phone: "+11234123412",
	})
	assert.Equal(t, response.Status, proto.Status_ERROR)
	assert.Equal(t, *response.Detail, "User exists")
}

func TestSignupSuccessFlow(t *testing.T) {
	authService := &services.AuthService{
		UserRepository: &mocks.UserRepositoryMock{
			CreateResponse: true,
		},
		OTPRepository: &mocks.OTPRepositoryMock{
			CreateResponse: true,
		},
		OTP: &mocks.OTPMock{
			SendOTPResponse: true,
		},
	}
	response, _ := authService.SignupWithPhoneNumber(context.Background(), &proto.SignupRequest{
		Name:  "Sergio",
		Phone: "+11234123412",
	})
	assert.Equal(t, response.Status, proto.Status_SUCCESS)
}
