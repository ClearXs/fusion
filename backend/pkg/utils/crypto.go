package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/samber/lo"
	"strings"
)

// EncryptPassword 给定username、password和salt，采取sha256加密
func EncryptPassword(username string, password string, salt string) string {
	if lo.IsEmpty(username) || lo.IsEmpty(password) || lo.IsEmpty(salt) {
		return ""
	}
	return Encrypt(Encrypt(username+Encrypt(password+salt)) + salt + Encrypt(username+salt))
}

// Encrypt 给定文本采取sha256加密
func Encrypt(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

// MakeSalt 基于random string与base64创建随机的slat value
func MakeSalt() string {
	random := lo.RandomString(32, lo.LettersCharset)
	bb := &strings.Builder{}
	encoder := base64.NewEncoder(base64.StdEncoding, bb)
	_, err := encoder.Write([]byte(random))
	if err != nil {
		return ""
	}
	return bb.String()
}
