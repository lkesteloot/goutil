package webutil

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"strconv"
	"strings"
	"time"
)

// Creates a string that can be used to sign a secure cookie.
func makeCookieSignature(name, value, timestamp string, signedCookieKey []byte) string {
	hash := hmac.New(sha1.New, signedCookieKey)
	// Use a separator here so that attacker can't transfer trailing digits from
	// value to timestamp.
	hash.Write(([]byte)(name + "|" + value + "|" + timestamp))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

// Returns the string that should be placed in a cookie value. The cookie is
// not encrypted, but it is signed and timestamped to prevent tampering.
func MakeSecureCookie(name string, value, signedCookieKey []byte) string {
	value64 := base64.URLEncoding.EncodeToString(value)
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	signature := makeCookieSignature(name, value64, timestamp, signedCookieKey)

	return value64 + "|" + timestamp + "|" + signature
}

// Takes a cookie value as created by MakeSecureCookie() and returns the
// payload.
func ParseSecureCookie(name, cookieValue string, signedCookieKey []byte,
	signatureExpiration time.Duration) (string, bool) {

	fields := strings.Split(cookieValue, "|")
	if len(fields) != 3 {
		log.Printf("Invalid cookie format: %s", cookieValue)
		return "", false
	}

	signature := makeCookieSignature(name, fields[0], fields[1], signedCookieKey)
	if signature != fields[2] {
		log.Printf("Invalid cookie signature: %s", cookieValue)
		return "", false
	}

	timestamp, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		log.Printf("Invalid cookie timestamp: %s", cookieValue)
		return "", false
	}

	if signatureExpiration != 0 {
		if timestamp < time.Now().UTC().Add(-signatureExpiration).Unix() {
			log.Printf("Expired cookie: %s", cookieValue)
			return "", false
		}
	}

	payload, err := base64.URLEncoding.DecodeString(fields[0])
	if err != nil {
		log.Printf("Can't decode payload: %s", cookieValue)
		return "", false
	}

	return string(payload), true
}
