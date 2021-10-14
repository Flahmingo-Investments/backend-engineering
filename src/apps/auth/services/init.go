package services

import (
	jwt "grpc-example/src/commons/services/jwt"
	"grpc-example/src/commons/services/pubsub"
	"grpc-example/src/data/repositories"
)

type AuthService struct {
	UserRepository repositories.UserManager
	OTPRepository  repositories.OTPManager
	OTP            pubsub.OTPPubSub
	JWT            jwt.JWTManager
}

func Initialize() *AuthService {
	return &AuthService{
		UserRepository: &repositories.UserRepository{},
		OTPRepository:  &repositories.OTPRepository{},
		OTP:            &pubsub.OTPPublisher{},
		JWT:            jwt.Initialize(),
	}
}
