package services

import (
	"context"
	commons "grpc-example/src/commons/utils"
	"grpc-example/src/data/models"
	"grpc-example/src/proto"
)

func validateSignupData(message *proto.SignupRequest) *proto.GenericResponse {
	var validation string
	if commons.IsEmpty(message.Name) {
		validation = "Empty name"
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

func validateUserExists(service *AuthService, message *proto.SignupRequest) *proto.GenericResponse {
	user := service.UserRepository.Get(&models.User{
		Phone: message.Phone,
	})

	if user != nil {
		validation := "User exists"
		return &proto.GenericResponse{
			Status: proto.Status_ERROR,
			Detail: &validation,
		}
	}
	return nil
}

func (service *AuthService) SignupWithPhoneNumber(ctx context.Context, message *proto.SignupRequest) (*proto.GenericResponse, error) {
	result := validateSignupData(message)
	if result != nil {
		return result, nil
	}

	result = validateUserExists(service, message)
	if result != nil {
		return result, nil
	}

	user := &models.User{
		Name:  message.Name,
		Phone: message.Phone,
	}
	saved := service.UserRepository.Create(user)
	if !saved {
		validation := "Error creating user on database"
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
