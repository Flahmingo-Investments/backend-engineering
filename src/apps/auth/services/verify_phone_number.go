package services

import (
	"context"
	commons "grpc-example/src/commons/utils"
	"grpc-example/src/proto"
)

func validateVerifyPhoneData(message *proto.OTPValidationRequest) *proto.GenericResponse {
	var validation string
	if commons.IsEmpty(message.Otp) {
		validation = "Empty otp"
	} else if commons.IsEmpty(message.Phone) {
		validation = "Empty phone"
	} else if !commons.IsPhoneNumber(message.Phone) {
		validation = "Phone number structure: +1012345789"
	}
	if !commons.IsEmpty(validation) {
		return &proto.GenericResponse{
			Status: proto.Status_ERROR,
			Detail: &validation,
		}
	}
	return nil
}

func (service *AuthService) VerifyPhoneNumber(ctx context.Context, message *proto.OTPValidationRequest) (*proto.GenericResponse, error) {
	result := validateVerifyPhoneData(message)
	if result != nil {
		return result, nil
	}
	otp := service.OTPRepository.GetByPhone(message)
	if otp == nil {
		validation := "Otp doesn't exist"
		return &proto.GenericResponse{
			Status: proto.Status_ERROR,
			Detail: &validation,
		}, nil
	}

	otp.Used = true
	service.OTPRepository.Update(otp)

	otp.User.Verified = true
	service.UserRepository.Update(&otp.User)

	return &proto.GenericResponse{
		Status: proto.Status_SUCCESS,
	}, nil
}
