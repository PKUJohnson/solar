package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"math/big"
	"time"

	std "github.com/PKUJohnson/solar/std"
	"golang.org/x/crypto/bcrypt"
)

const (
	numberchars = "1234567890"
)

func GetMD5Content(content string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(content))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func GetMD5Content2(content []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(content)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func Base64Decode(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

func Base64URLEncode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func Base64URLDecode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}

func RandInt63(max int64) int64 {
	var maxbi big.Int
	maxbi.SetInt64(max)
	value, _ := rand.Int(rand.Reader, &maxbi)
	return value.Int64()
}

func RandNumStr(l int) string {
	ret := make([]byte, 0, l)
	for i := 0; i < l; i++ {
		index := RandInt63(int64(len(numberchars)))
		ret = append(ret, numberchars[index])
	}
	return string(ret)
}

func RandBytes(c int) ([]byte, error) {
	b := make([]byte, c)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

func SimpleGuid() string {
	b := make([]byte, 24)
	binary.LittleEndian.PutUint64(b, uint64(time.Now().UnixNano()))
	if _, err := rand.Read(b[8:]); err != nil {
		return ""
	} else {
		return base64.RawURLEncoding.EncodeToString(b)
	}
}

func SaltPassword(password string) []byte {
	h, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		std.LogErrorLn("Bcrypt salt password error: ", err)
		return nil
	}
	return h
}

func VerifyPassword(stored []byte, password string) bool {
	return bcrypt.CompareHashAndPassword(stored, []byte(password)) == nil
}

func SaltPassword2(password string) []byte {
	if rb, err := RandBytes(64); err != nil {
		return nil
	} else {
		pw := append(rb, []byte(password)...)
		hash := sha512.Sum512(pw)
		return append(rb, hash[:]...)
	}
}

func VerifyPassword2(stored []byte, password string) bool {
	if len(stored) != 128 {
		return false
	}
	pw := append([]byte{}, stored[:64]...)
	pw = append(pw, []byte(password)...)
	hash := sha512.Sum512(pw)
	return bytes.Equal(hash[:], stored[64:])
}

func Md5(data []byte) []byte {
	sum := md5.Sum(data)
	return sum[:]
}

func Md5Str(data []byte) string {
	return base64.StdEncoding.EncodeToString(Md5(data))
}

func Encrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := padding(data)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return appendHmac(ciphertext, key), nil
}

func Decrypt(data, key []byte) ([]byte, error) {
	data, err := removeHmac(data, key)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(data)%aes.BlockSize != 0 {
		errors.New("data is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(data))
	mode.CryptBlocks(plaintext, data)
	return unPadding(plaintext), nil
}

func padding(data []byte) []byte {
	blockSize := aes.BlockSize
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func unPadding(data []byte) []byte {
	l := len(data)
	unpadding := int(data[l-1])
	return data[:(l - unpadding)]
}

func appendHmac(data, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	hash := mac.Sum(nil)
	return append(data, hash...)
}

func removeHmac(data, key []byte) ([]byte, error) {
	if len(data) < sha256.Size {
		return nil, errors.New("Invalid length")
	}
	p := len(data) - sha256.Size
	mmac := data[p:]
	mac := hmac.New(sha256.New, key)
	mac.Write(data[:p])
	exp := mac.Sum(nil)
	if hmac.Equal(mmac, exp) {
		return data[:p], nil
	} else {
		return nil, errors.New("MAC doesn't match")
	}
}
