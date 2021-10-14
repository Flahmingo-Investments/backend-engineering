package mocks

import (
	"grpc-example/src/data/models"
	"grpc-example/src/proto"
)

type OTPRepositoryMock struct {
	CreateResponse     bool
	GetByPhoneResponse *models.OTP
	UpdateResponse     bool
	UpdateOTPCache     *models.OTP
}

func (mock *OTPRepositoryMock) Create(*models.OTP) bool {
	return mock.CreateResponse
}

func (mock *OTPRepositoryMock) Update(otp *models.OTP) bool {
	mock.UpdateOTPCache = otp
	return mock.UpdateResponse
}

func (mock *OTPRepositoryMock) GetByPhone(*proto.OTPValidationRequest) *models.OTP {
	return mock.GetByPhoneResponse
}
