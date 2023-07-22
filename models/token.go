package models

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Token struct {
	jwt.StandardClaims
}

type AccessToken struct {
	UserID         string `json:"uid"`
	RefreshTokenID string `json:"rt_id"`
	Token
}

func (token *AccessToken) SetExpiration(expiryTime time.Time) *AccessToken {
	token.ExpiresAt = expiryTime.Unix()
	return token
}

func (token *AccessToken) ToJWT() *jwt.Token {
	return jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)
}

func (token *AccessToken) ToJWTString() string {
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)
	tokenString, _ := jwtToken.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	return tokenString
}

type RefreshToken struct {
	UserID string `json:"uid"`
	Token
}

func (refreshToken *RefreshToken) ToJWT() *jwt.Token {
	return jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshToken)
}

func (refreshToken *RefreshToken) ToJWTString() string {
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshToken)
	tokenString, _ := jwtToken.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	return tokenString
}

type TokenSet struct {
	ID                 string `gorm:"primary_key"`
	UserID             string
	RefreshTokenID     string  `gorm:"column:rt_id;unique"`
	PrevRefreshTokenID *string `gorm:"column:prt_id;unique"`
	UpdatedAt          time.Time
}
