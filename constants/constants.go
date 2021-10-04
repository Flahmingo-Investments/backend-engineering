package constants

const PUBSUB_SENDOTP_MSG_STRING_SEPERATOR = "---"
const VERIFICATION_INSTR_MSG_STRING = "OTP for Phone Verification "
const LOGIN_INSTR_MSG_STRING = "OTP for Login "

// Configurations
const OTP_MAX_INT = 8999
const OTP_MIN_INT = 1000
const PUBSUB_PROJECT_ID = "flahmingo-327520"
const PUBSUB_TOPIC = "sendOTP" // Change to "verification"
const PUBSUB_SUBSCRIPTION_ID = "sendOTPServer"
