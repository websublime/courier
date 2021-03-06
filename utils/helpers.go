package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/dgrijalva/jwt-go"
)

const (
	ErrorServerUnknown           = "ESERVER_UNKNOWN"
	ErrorBodyParse               = "EBODY_PARSER"
	ErrorStatusForbidden         = "ESTATUS_FORBIDDEN"
	ErrorInvalidToken            = "ETOKEN_INVALID"
	ErrorEncryptFailure          = "ECRYPT_FAILURE"
	ErrorAudienceCreate          = "EAUD_CREATEFAILURE"
	ErrorAudienceNotFound        = "EAUD_NOTFOUND"
	ErrorChannelCreation         = "ECHANNEL_CREATEFAILURE"
	ErrorHookInvalid             = "EHOOK_INVALID"
	ErrorTopicsEmpty             = "ETOPICS_EMPTY"
	ErrorTopicsInvalid           = "ETOPICS_INVALID"
	ErrorTopicUnknown            = "ETOPICS_UNKNOWN"
	ErrorTopicNotFound           = "ETOPICS_NOTFOUND"
	ErrorTopicCreate             = "ETOPICS_CREATE"
	ErrorSocketConnectionFailure = "ESOCKET_CONNECTIONFAILURE"
	ErrorSocketMessageFailure    = "ESOCKET_MESSAGEFAILURE"
)

// Contains checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

func Decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

func ParseJwtToken(tokenString string, secret string) (*jwt.Token, error) {
	parser := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}

	token, err := parser.ParseWithClaims(tokenString, &GoTrueClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
