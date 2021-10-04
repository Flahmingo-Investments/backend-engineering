package cryptUtils

/** Idea -> To serve as utils class to store all commonly used crypt functions and implementations
To do:
1)
2)
*/

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTPMsg(min, max int64) string {
	if (min < 0) || (max < 0) {
		panic("Invalid Range")
	}
	var otpInt int64
	otpBigInt, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		panic(err)
	}
	if otpBigInt.IsInt64() {
		otpInt = otpBigInt.Int64() + min
	}
	return fmt.Sprint(otpInt)
}
