package mocks

import "grpc-example/src/proto"

type OTPMock struct {
	SendOTPResponse bool
}

func (mock *OTPMock) SendOTP(*proto.OTPValidationRequest) bool {
	return mock.SendOTPResponse
}
