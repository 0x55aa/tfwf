package tfwf

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func inSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func get_random_string(length int, char []byte) (s string) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		s = s + string(char[rand.Intn(len(char))])
	}
	return
}

func cookie_signature(name, value, timestamp string) string {
	//这里取了settings
	hash := hmac.New(sha256.New, []byte(Settings["secret_key"]))
	hash.Write([]byte(name + value + timestamp))

	return hex.EncodeToString(hash.Sum(nil))
}

func create_signed_value(name string, value string) string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	value = base64.StdEncoding.EncodeToString([]byte(value))
	signature := cookie_signature(name, value, timestamp)
	value = strings.Join([]string{value, timestamp, signature}, "|")

	return value
}
