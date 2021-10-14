package services

import (
	"context"
	"encoding/base64"
	"grpc-example/src/proto"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func (service *OTPService) SendOTP(ctx context.Context, message *proto.OTPValidationRequest) (*proto.GenericResponse, error) {
	data := url.Values{}
	data.Set("From", os.Getenv("TWILIO_PHONE_NUMBER"))
	data.Set("Body", message.Otp)
	data.Set("To", message.Phone)

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	token := os.Getenv("TWILIO_AUTH_TOKEN")
	endpoint := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	client := &http.Client{}
	request, _ := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	basicAuth := base64.StdEncoding.EncodeToString([]byte(accountSid + ":" + token))
	request.Header.Add("Authorization", "Basic "+basicAuth)

	response, err := client.Do(request)

	var result bool
	if err != nil {
		log.Printf("Error %s", err.Error())
		result = false
	} else {
		result = response.StatusCode == 201
		if !result {
			body, _ := ioutil.ReadAll(response.Body)
			log.Print(string(body))
		}
	}
	if result {
		return &proto.GenericResponse{
			Status: proto.Status_SUCCESS,
		}, nil
	}
	validation := "Error in twilio"
	return &proto.GenericResponse{
		Status: proto.Status_ERROR,
		Detail: &validation,
	}, nil
}
