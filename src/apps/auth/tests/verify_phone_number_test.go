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

func TestVerifyPhoneDataValidation(t *testing.T) {
	authService := services.AuthService{}
	response, _ := authService.VerifyPhoneNumber(context.Background(), &proto.OTPValidationRequest{
		Otp:   "",
		Phone: "0123456789",
	})
	assert.Equal(t, response.Status, proto.Status_ERROR)
	assert.Equal(t, *response.Detail, "Empty otp")

	response, _ = authService.VerifyPhoneNumber(context.Background(), &proto.OTPValidationRequest{
		Otp:   "123",
		Phone: "",
	})
	assert.Equal(t, response.Status, proto.Status_ERROR)
	assert.Equal(t, *response.Detail, "Empty phone")

	response, _ = authService.VerifyPhoneNumber(context.Background(), &proto.OTPValidationRequest{
		Otp:   "1234",
		Phone: "1234",
	})
	assert.Equal(t, response.Status, proto.Status_ERROR)
	assert.Equal(t, *response.Detail, "Phone number structure: +1012345789")
}

func TestVerifyPhoneOTPDoesNotExist(t *testing.T) {
	authService := &services.AuthService{
		OTPRepository: &mocks.OTPRepositoryMock{},
	}
	response, _ := authService.VerifyPhoneNumber(context.Background(), &proto.OTPValidationRequest{
		Otp:   "1234",
		Phone: "+11234123412",
	})
	assert.Equal(t, response.Status, proto.Status_ERROR)
	assert.Equal(t, *response.Detail, "Otp doesn't exist")
}

func TestVerifyPhoneSuccessFlow(t *testing.T) {
	otpMock := &mocks.OTPRepositoryMock{
		GetByPhoneResponse: &models.OTP{
			Code: "1234",
		},
	}
	userMock := &mocks.UserRepositoryMock{
		UpdateResponse: true,
	}
	authService := &services.AuthService{
		OTPRepository:  otpMock,
		UserRepository: userMock,
	}
	response, _ := authService.VerifyPhoneNumber(context.Background(), &proto.OTPValidationRequest{
		Otp:   "1234",
		Phone: "+11234123412",
	})
	assert.Equal(t, response.Status, proto.Status_SUCCESS)
	assert.Equal(t, true, otpMock.UpdateOTPCache.Used)
	assert.Equal(t, true, otpMock.UpdateOTPCache.User.Verified)
}
