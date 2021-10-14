package services

import (
	"context"
	"grpc-example/src/data/models"
	"grpc-example/src/proto"
	"time"

	"github.com/golang-jwt/jwt"
)

func (service *AuthService) getLoggedUser(message *proto.ProfileRequest) *models.User {
	user := &models.User{}
	user.ID = service.JWT.ParseJWT(message.AccessToken, jwt.SigningMethodHS256)
	if user.ID == "" {
		return nil
	}
	user = service.UserRepository.Get(user)
	return user
}

func (service *AuthService) GetProfile(ctx context.Context, message *proto.ProfileRequest) (*proto.ProfileResponse, error) {
	user := service.getLoggedUser(message)
	if user == nil {
		return &proto.ProfileResponse{
			Status: proto.Status_UNAUTHORIZED,
		}, nil
	}
	return &proto.ProfileResponse{
		Status:    proto.Status_SUCCESS,
		Id:        user.ID,
		Name:      user.Name,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.Format(time.ANSIC),
		UpdatedAt: user.UpdatedAt.Format(time.ANSIC),
		Verified:  user.Verified,
		Logged:    user.Logged,
	}, nil
}
