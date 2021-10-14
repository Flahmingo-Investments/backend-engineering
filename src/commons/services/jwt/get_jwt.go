package jwt

import (
	commons "grpc-example/src/commons/utils"
	"grpc-example/src/data/models"
	"time"

	jwtLib "github.com/golang-jwt/jwt"
)

func (service *JWTAuth) GetJWT(user *models.User, method jwtLib.SigningMethod, duration time.Duration) string {
	atClaims := jwtLib.MapClaims{}
	atClaims["user_id"] = user.ID
	atClaims["exp"] = commons.Now().Add(duration).Unix()
	at := jwtLib.NewWithClaims(method, atClaims)
	token, err := at.SignedString([]byte(service.SECRET))
	if err != nil {
		return ""
	}
	return token
}
