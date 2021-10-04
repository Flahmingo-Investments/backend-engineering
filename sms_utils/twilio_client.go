package twilio

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendOTP(phoneNumberTo, message string) {

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_ACCOUNT_AUTH_TOKEN")
	phoneNumberFrom := os.Getenv("TWILIO_FROM_PHONE_NUMBER")
	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(phoneNumberTo)
	params.SetFrom(phoneNumberFrom)
	params.SetBody(message)

	resp, err := client.ApiV2010.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}
}
