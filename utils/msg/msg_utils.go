package message

import (
	constants "flahmingo/constants"
	"fmt"
	"strings"
)

/** Idea -> To serve as a template for messages sent on pubsub / msgbroker
1) Develop a standard format for message to be sent and received through pubsub
2) Extract phoneNumber from the encoded msg string
3) Extract otpCode from the encoded msg string
4) Extract instruction from the encoded msg string
*/

type OTPMsg struct {
	PhoneNumber, OtpCode, InstructionMessage string
	// To do : timestamp for auto drop
}

func CreateMsg(phoneNumber, otpCode, message string) string {
	return OTPMsg{PhoneNumber: phoneNumber, OtpCode: otpCode, InstructionMessage: message}.toString()
}

func (msg OTPMsg) toString() string {
	return fmt.Sprintf("%s%s%s%s%s", msg.PhoneNumber, constants.PUBSUB_SENDOTP_MSG_STRING_SEPERATOR, msg.OtpCode, constants.PUBSUB_SENDOTP_MSG_STRING_SEPERATOR, msg.InstructionMessage)
}
func fromString(messageString string) OTPMsg {
	var s []string
	s = strings.SplitN(messageString, constants.PUBSUB_SENDOTP_MSG_STRING_SEPERATOR, 3)
	return OTPMsg{PhoneNumber: s[0], OtpCode: s[1], InstructionMessage: s[2]}
}

func GetPhoneNumber(msg string) string {
	return fromString(msg).PhoneNumber
}

func GetInstructionMessage(msg string) string {
	return fromString(msg).InstructionMessage
}
func GetOtpCode(msg string) string {
	return fromString(msg).OtpCode
}

