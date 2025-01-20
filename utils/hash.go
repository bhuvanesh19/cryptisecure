package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"math"
	"strconv"
	"time"
)

const TIME_STEP = 30
const T0 = 0
const HOTP_LENGTH = 6

func GenerateTOTP(K string) string {

	T := (time.Now().Unix() - T0) / TIME_STEP
	return GenerateHOTP(uint32(T), K)
}

func GenerateHOTP(C uint32, K string) string {
	hasher := hmac.New(sha1.New, []byte(K))
	c_slice := make([]byte, 8)

	for i := range c_slice {
		c_slice[7-i] = byte(C & 0xff)
		C = C >> 8
	}
	_, err := hasher.Write(c_slice)
	if err != nil {
		panic(err)
	}
	hashed_value := hasher.Sum(nil)
	return dynamicTruncation(HOTP_LENGTH, hashed_value)

}
func dynamicTruncation(digits uint, data []byte) string {
	offset := int(data[19] & 0x0f)
	bin_code := make([]byte, 4)
	bin_code[0] = data[offset] & 0x7f
	bin_code[1] = data[offset+1]
	bin_code[2] = data[offset+2]
	bin_code[3] = data[offset+3]
	otp := strconv.Itoa(int(stToNum(bin_code) % uint64(math.Pow10(int(digits)))))
	if len(otp) < int(digits) {
		for i := 0; i < int(digits)-len(otp); i++ {
			otp = "0" + otp
		}
	}
	return otp

}
func stToNum(bits []byte) uint64 {
	var counter int = len(bits) - 1
	var ret uint64 = 0
	for _, i := range bits {
		ret += (uint64(math.Pow(2, (float64(counter)*8))) * uint64(int(i)))
		counter -= 1
	}

	return ret
}

func CalculateSHA256(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}
