package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"gitlab.com/vita/go/verify"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func requestOTP(phoneNumber string) {
	accountSid := "REDACTED"
	authToken := "REDACTED"

	client := verify.NewClient(accountSid, authToken)

	number := "+1" + phoneNumber

	fmt.Print(number)

	_, err := client.NewVerification(
		context.TODO(), "REDACTED",
		&verify.VerificationInput{
			To:      number,
			Channel: verify.ChannelSMS,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

}
func verifyOTP(phoneNumber string, otp string) (valid bool) {
	accountSid := "REDACTED"
	authToken := "REDACTED"
	number := "+1" + phoneNumber
	fmt.Print(number)
	valid = false

	client := verify.NewClient(accountSid, authToken)
	out, err := client.NewVerificationCheck(
		context.TODO(), "REDACTED",
		&verify.VerificationCheckInput{
			To:   number,
			Code: otp,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if out != nil {
		switch out.Status {
		case verify.StatusApproved:
			fmt.Println("Approved!")
			valid = true
		}
	}
	return valid
}

func signupWithPhoneNumber(phoneNumber string) {
	ctx := context.Background()

	projectID := "pubsubotp" //gcp project ID

	// Initialize client
	auth_service, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("../keys/auth_key.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer auth_service.Close()

	topic := auth_service.Topic("verification")

	var results []*pubsub.PublishResult
	res := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(phoneNumber + " " + "sendOTP"), //publish phone number and sendOTP request
	})
	results = append(results, res)

	topic.Stop()
}

//get phone number from session cookies
func getPhoneNumber(request *http.Request) (phoneNumber string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			phoneNumber = cookieValue["phoneNumber"]
		}
	}
	return phoneNumber
}

func setSession(phoneNumber string, fname string, lname string, dob string, response http.ResponseWriter) {
	value := map[string]string{
		"phoneNumber": phoneNumber,
		"fname":       fname,
		"lname":       lname,
		"dob":         dob,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// login handler

func loginHandler(response http.ResponseWriter, request *http.Request) {
	phoneNumber := request.FormValue("phn")
	fname := request.FormValue("fname")
	lname := request.FormValue("lname")
	dob := request.FormValue("dob")
	redirectTarget := "/"

	fmt.Println(phoneNumber)

	setSession(phoneNumber, fname, lname, dob, response)

	fmt.Println(phoneNumber)
	redirectTarget = "/verify"

	http.Redirect(response, request, redirectTarget, 302)
	signupWithPhoneNumber(phoneNumber)
	requestOTP(phoneNumber)

}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// index page

const indexPage = `
	<h1>Welcome to Flahmingo!</h1>
	<h2>Enter your phone number to get started:</h2>
	<form method="post" action="/login">
		<label for="phn">Phone number: </label>
		<input type="tel" id="phn" name="phn" pattern="[0-9]{10}" placeholder="1234567890"
        required>
		<br>
		<label for="fname">First Name: </label>
		<input type="text" id="fname" name="fname" required>
		<br>
		<label for="lname">Last Name: </label>
		<input type="text" id="lname" name="lname" required>
		<br>
		<label for="dob">Date of Birth: </label>
		<input type="date" id="dob" name="dob" required>
		<br>
		<button type="submit">Create Account</button>
	</form>
	`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

// profile page

const profilePage = `
	<h1>Profile</h1>
	<hr>
	<small>User: %s</small>
	<form method="post" action="/logout">
		<button type="submit">Logout</button>
	</form>
	`

func profilePageHandler(response http.ResponseWriter, request *http.Request) {
	phoneNumber := getPhoneNumber(request)

	otp := request.FormValue("otp")

	fmt.Println(otp)

	if otp != "" {
		if verifyOTP(getPhoneNumber(request), otp) {
			fmt.Fprintf(response, profilePage, phoneNumber)
		}
	} else {
		log.Fatal("unable to verify")

	}
}

// verification page

const verificationPage = `
<h1>Verification</h1>
<hr>
<h2>Not so fast...</h2>
<p>Please enter the OTP sent to %s.</p>
<form method="post" action="/profile">
	<input type="number" id="otp" name="otp" required>
	<button type="submit">Verify</button>
</form>
`

func verificationPageHandler(response http.ResponseWriter, request *http.Request) {

	fmt.Fprintf(response, verificationPage, getPhoneNumber(request))

}

// Instantiate Gorilla/mux server for HTTP routing
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// server main method

var router = mux.NewRouter()

func main() {
	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/profile", profilePageHandler)
	router.HandleFunc("/verify", verificationPageHandler)
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
