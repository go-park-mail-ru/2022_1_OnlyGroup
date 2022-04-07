package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtToken struct {
	Secret []byte
}

func NewJwtToken(secret string) (*JwtToken, error) {
	return &JwtToken{Secret: []byte(secret)}, nil
}

type JwtCsrfClaims struct {
	Session string `json:"sid"`
	UserID  int    `json:"uid"`
	jwt.StandardClaims
}

func (tk *JwtToken) Create(session string, id int, tokenExpTime int64) (string, error) {
	data := JwtCsrfClaims{
		Session: session,
		UserID:  id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpTime,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString(tk.Secret)
}

func (tk *JwtToken) parseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, jwt.ErrInvalidKey
	}
	return tk.Secret, nil
}

func (tk *JwtToken) Check(session string, id int, inputToken string) error {
	payload := &JwtCsrfClaims{}
	_, err := jwt.ParseWithClaims(inputToken, payload, tk.parseSecretGetter)
	if err != nil {
		return ErrBadCSRF.Wrap(err, "")
	}
	if payload.Valid() != nil {
		return ErrBadCSRF.Wrap(jwt.ValidationError{Inner: jwt.ErrInvalidKey, Errors: jwt.ValidationErrorExpired}, "")
	}
	if payload.Session != session && payload.UserID != id {
		return ErrBaseApp.Wrap(jwt.ValidationError{Inner: jwt.ErrInvalidKey, Errors: jwt.ValidationErrorId}, "")

	}

	return nil
}
