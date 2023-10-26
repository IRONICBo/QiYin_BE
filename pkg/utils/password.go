package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

const (
	_salt = "QiYin"
)

// EncryptPassword encrypt password.
func EncryptPassword(password string) string {
	// md5
	m := md5.New()
	m.Write([]byte(password + _salt))
	mByte := m.Sum(nil)

	// hmac
	h := hmac.New(sha256.New, []byte(_salt))
	h.Write(mByte)
	password = hex.EncodeToString(h.Sum(nil))

	return password
}

// ComparePassword compare password.
func ComparePassword(password, encryptedPassword string) bool {
	return EncryptPassword(password) == encryptedPassword
}
