package services

import (
	"context"
	commons "grpc-example/src/commons/utils"
	"grpc-example/src/proto"
	"time"

	"github.com/golang-jwt/jwt"
)

func validateLoginPhoneData(message *proto.OTPValidationRequest) *proto.JWTResponse {
	var validation string
	if commons.IsEmpty(message.Otp) {
		validation = "Empty otp"
	} else if commons.IsEmpty(message.Phone) {
		validation = "Empty phone"
	} else if !commons.IsPhoneNumber(message.Phone) {
		validation = "Phone number structure: +1012345789"
	}
	if !commons.IsEmpty(validation) {
		return &proto.JWTResponse{
			Status: proto.Status_ERROR,
			Detail: &validation,
		}
	}
	return nil
}

func (service *AuthService) ValidatePhoneNumberLogin(ctx context.Context, message *proto.OTPValidationRequest) (*proto.JWTResponse, error) {
	result := validateLoginPhoneData(message)
	if result != nil {
		return result, nil
	}
	otp := service.OTPRepository.GetByPhone(message)
	if otp == nil {
		validation := "Otp doesn't exist"
		return &proto.JWTResponse{
			Status: proto.Status_ERROR,
			Detail: &validation,
		}, nil
	}

	otp.Used = true
	service.OTPRepository.Update(otp)

	otp.User.Logged = true
	service.UserRepository.Update(&otp.User)

	accessToken := service.JWT.GetJWT(&otp.User, jwt.SigningMethodHS256, time.Minute*10)
	refreshToken := service.JWT.GetJWT(&otp.User, jwt.SigningMethodHS512, time.Hour*24*7)

	return &proto.JWTResponse{
		Status:       proto.Status_SUCCESS,
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
	}, nil
}
