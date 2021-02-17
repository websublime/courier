package config

import "github.com/dgrijalva/jwt-go"

type GoTrueClaims struct {
	jwt.StandardClaims
	Email        string                 `json:"email"`
	AppMetaData  map[string]interface{} `json:"app_metadata"`
	UserMetaData map[string]interface{} `json:"user_metadata"`
}
