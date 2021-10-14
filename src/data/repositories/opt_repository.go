package repositories

import (
	commons "grpc-example/src/commons/utils"
	"grpc-example/src/data/models"
	"grpc-example/src/proto"

	uuid "github.com/satori/go.uuid"
)

type OTPManager interface {
	Create(*models.OTP) bool
	Update(*models.OTP) bool
	GetByPhone(*proto.OTPValidationRequest) *models.OTP
}

type OTPRepository struct {
}

func (u *OTPRepository) Create(otp *models.OTP) bool {
	otp.ID = uuid.NewV4().String()
	result := commons.DB.Create(otp)
	return result.Error == nil
}

func (u *OTPRepository) Update(otp *models.OTP) bool {
	result := commons.DB.Save(otp)
	return result.Error == nil && result.RowsAffected > 0
}

func (u *OTPRepository) GetByPhone(search *proto.OTPValidationRequest) *models.OTP {
	var otp *models.OTP
	result := commons.DB.Joins("User", commons.DB.Where(&models.User{Phone: search.Phone})).Where(&models.OTP{
		Code: search.Otp,
		Used: false,
	}).Find(&otp)
	if result.RowsAffected > 0 {
		return otp
	}
	return nil
}
