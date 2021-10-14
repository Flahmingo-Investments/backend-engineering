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

func TestLoginDataValidation(t *testing.T) {
	authService := services.AuthService{}
	response, _ := authService.LoginWithPhoneNumber(context.Background(), &proto.PhoneRequest{
		Phone: "0123456789",
	})
	assert.Equal(t, response.Status, proto.Status_ERROR)
	assert.Equal(t, *response.Detail, "Phone number structure: +1012345789")

	response, _ = authService.LoginWithPhoneNumber(context.Background(), &proto.PhoneRequest{
		Phone: "",
	})
	assert.Equal(t, response.Status, proto.Status_ERROR)
	assert.Equal(t, *response.Detail, "Empty phone")
}
func TestLoginPhoneNotFound(t *testing.T) {
	authService := &services.AuthService{
		UserRepository: &mocks.UserRepositoryMock{},
	}
	response, _ := authService.LoginWithPhoneNumber(context.Background(), &proto.PhoneRequest{
		Phone: "+11234123412",
	})
	assert.Equal(t, response.Status, proto.Status_ERROR)
	assert.Equal(t, *response.Detail, "User doesn't exist")
}

func TestLogingSuccessFlow(t *testing.T) {
	authService := &services.AuthService{
		UserRepository: &mocks.UserRepositoryMock{
			GetResponse: &models.User{},
		},
		OTPRepository: &mocks.OTPRepositoryMock{
			CreateResponse: true,
		},
		OTP: &mocks.OTPMock{
			SendOTPResponse: true,
		},
	}
	response, _ := authService.LoginWithPhoneNumber(context.Background(), &proto.PhoneRequest{
		Phone: "+11234123412",
	})
	assert.Equal(t, response.Status, proto.Status_SUCCESS)
}
