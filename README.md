# Home work

## What is it?

In this exercise, you'll build two microservices that will contact over Google Cloud Pub/Sub. Below is a list of requirements and success criteria for your finished project.

## Requirements

You'll build two `gRPC` services, `auth-service` and `otp-service`. `auth-service` will publish message `SendOTP` on `verification` topic, and the `otp-service` will subscribe to the topic and will listen on the event.
After consuming the event `otp-service` will use Twilio to send the otp.

You'll build this service in `golang` and use GCP Cloud Pub/Sub for messaging.


### Auth Service

Create an auth service. This service will handle authentication and user profile.

#### Create Account

1. User should be able to create account using phone number.
2. On successful account creation, publish a message on pubsub topic.
3. User should be able to verify the account using the received otp.
4. User should be able to login. By using otp.
5. User should be able to logout.
6. User should be able to get their profile.

#### What to implement
1. SignupWithPhoneNumber
2. VerifyPhoneNumber
3. LoginWithPhoneNumber
4. ValidatePhoneNumberLogin
5. GetProfile

### OTP Service

1. Will consume `SendOTP` event.
2. Send OTP using Twilio.

#### What to implement
1. SendOTP

# Resources

## GCP Pub/Sub

https://cloud.google.com/pubsub/docs/overview

## Twilio

https://www.twilio.com/docs/sms/api


# Bonus Points

1. Tests
2. Comments

# How to contribute 

Contributing
Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are greatly appreciated.

Fork the project
Create your feature branch (git checkout -b feature/AmazingFeature)
Make your changes
Commit your changes (git commit -m 'Add some AmazingFeature')
Push to the branch (git push origin feature/AmazingFeature)
Open a pull request


## Additional Instructions


- Makesure to fork the repo and send a PR ( if unsure check this https://medium.com/@rishabhmittal200/contributing-guide-when-you-fork-a-repository-3b97657b01fb)
