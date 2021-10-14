package jwt

import (
	commons "grpc-example/src/commons/utils"
	"time"

	jwtLib "github.com/golang-jwt/jwt"
)

func (service *JWTAuth) ParseJWT(tokenString string, method jwtLib.SigningMethod) string {
	token, _ := jwtLib.Parse(tokenString, func(token *jwtLib.Token) (interface{}, error) {
		return []byte(service.SECRET), nil
	})

	if token == nil || token.Method != method {
		return ""
	}
	data := token.Claims.(jwtLib.MapClaims)
	exp := time.Unix(int64(data["exp"].(float64)), 0)
	if exp.Before(commons.Now()) {
		return ""
	}
	signingSecret, _ := token.SigningString()
	secret := []byte(service.SECRET)
	err := method.Verify(signingSecret, token.Signature, secret)
	if err != nil {
		return ""
	}

	return data["user_id"].(string)
}
