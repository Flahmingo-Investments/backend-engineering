# Home work

## How to run
Add the following environment variables under ~/.profile
Note the program needs these values to be in place for it to work
| KEY  | VALUE |
| ---  | ----- |
| PATH | $PATH:/usr/local/go/bin |
| TWILIO_ACCOUNT_SID | "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" |
| TWILIO_ACCOUNT_AUTH_TOKEN | "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" |
| TWILIO_FROM_PHONE_NUMBER | "+1DDDDDDDDDD" |
| GOOGLE_APPLICATION_CREDENTIALS | "path/to/service/account/withPubSubPublisherSubscriberAccess.json" |

### Load the variables
~~~
$ source ~/.profile
~~~

### Changing configuration details of google pub/sub
Go to constants/constants.go file
| KEY  | VALUE |
| ---  | ----- |
| PUBSUB_PROJECT_ID | ANY_VALID_PROJECT_ID_CREATED_IN_GCP |
| PUBSUB_TOPIC | ANY_VALID_TOPIC_CREATED_IN_GCP |
| PUBSUB_SUBSCRIPTION_ID | ANY_VALID_SUBSCRIPTION_ID_CREATED_FOR_OTP_SERVICE |

### Running the servers
~~~
$ go run auth_service_server/main.go 
$ go run otp_service_server/main.go 
~~~

### Running the clients
~~~
$ go run auth_service_client/main.go SignupWithPhoneNumber +1DDDDDDDDDD  Amar
$ go run auth_service_client/main.go VerifyPhoneNumber +1DDDDDDDDDD  
$ go run auth_service_client/main.go LoginWithPhoneNumber +1DDDDDDDDDD  
$ go run auth_service_client/main.go ValidatePhoneNumberLogin +1DDDDDDDDDD  
$ go run auth_service_client/main.go GetProfile +1DDDDDDDDDD  Amar
$ go run auth_service_client/main.go Logout +1DDDDDDDDDD  Amar

$ go run otp_service_client/main.go +1DDDDDDDDDD
~~~

### To regenerate grpc stubs after any change to proto file

~~~
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc_service/grpcservice.proto
~~~


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
In order to test end to end the way you handle code commits, pull request and documentation. Follow this instructions. 

1. Fork the project
2. Create your feature branch (git checkout -b feature/AmazingFeature)
3. Make your changes
4. Commit your changes (git commit -m 'Add some AmazingFeature')
5. Push to the branch (git push origin feature/AmazingFeature)
6. Open a pull request


## Additional Instructions


- Makesure to fork the repo and send a PR ( if unsure check this https://medium.com/@rishabhmittal200/contributing-guide-when-you-fork-a-repository-3b97657b01fb)
