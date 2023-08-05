package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
	"unsafe"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var src = rand.NewSource(time.Now().UnixNano())

func ToSHA256(toHash string) string {
	hash := sha256.New()
	hash.Write([]byte(toHash))
	return hex.EncodeToString(hash.Sum(nil))
}

func RandString(n int) string {
	bytes := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			bytes[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&bytes))
}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func RandFileName(prefix, suffix string) string {
	return prefix + RandString(8) + suffix
}

func DataResponse(data interface{}, metadata interface{}) map[string]interface{} {
	return map[string]interface{}{"data": data, "meta": metadata}
}
