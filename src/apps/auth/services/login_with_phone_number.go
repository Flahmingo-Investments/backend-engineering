package services

import (
	"context"
	commons "grpc-example/src/commons/utils"
	"grpc-example/src/data/models"
	"grpc-example/src/proto"
)

func validateLoginData(message *proto.PhoneRequest) *proto.GenericResponse {
	var validation string
	if commons.IsEmpty(message.Phone) {
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

func (service *AuthService) LoginWithPhoneNumber(ctx context.Context, message *proto.PhoneRequest) (*proto.GenericResponse, error) {
	result := validateLoginData(message)
	if result != nil {
		return result, nil
	}

	user := service.UserRepository.Get(&models.User{Phone: message.Phone})
	if user == nil {
		validation := "User doesn't exist"
		return &proto.GenericResponse{
			Status: proto.Status_ERROR,
			Detail: &validation,
		}, nil
	}
	otp := &models.OTP{
		UserID: user.ID,
		Code:   commons.RandOTP(),
	}
	if service.OTPRepository.Create(otp) {
		service.OTP.SendOTP(&proto.OTPValidationRequest{
			Phone: user.Phone,
			Otp:   otp.Code,
		})
	}

	return &proto.GenericResponse{
		Status: proto.Status_SUCCESS,
	}, nil
}
