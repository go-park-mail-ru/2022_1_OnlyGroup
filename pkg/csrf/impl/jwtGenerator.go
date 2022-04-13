package impl

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtToken struct {
	Secret       []byte
	tokenExpTime int64
}

func NewJwtTokenGenerator(secret string, tokenExpTime int64) *JwtToken {
	return &JwtToken{Secret: []byte(secret), tokenExpTime: tokenExpTime}
}

type JwtCsrfClaims struct {
	Session string
	UserID  int
	URL     string
	jwt.StandardClaims
}

func (tk *JwtToken) Create(session string, id int, url string) (string, error) {
	data := JwtCsrfClaims{
		Session: session,
		UserID:  id,
		URL:     url,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tk.tokenExpTime,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	jwtToken, err := token.SignedString(tk.Secret)
	if err != nil {
		return "", handlers.ErrBaseApp.Wrap(err, "generate failed")
	}
	return jwtToken, nil
}

func (tk *JwtToken) ParseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, handlers.ErrBadCSRF
	}
	return tk.Secret, nil
}

func (tk *JwtToken) Check(session string, id int, url string, inputToken string) error {
	payload := &JwtCsrfClaims{}
	_, err := jwt.ParseWithClaims(inputToken, payload, tk.ParseSecretGetter)
	if err != nil {
		return handlers.ErrBadCSRF
	}
	if payload.Valid() != nil {
		return handlers.ErrBadCSRF
	}
	if payload.Session != session || payload.UserID != id || payload.URL != url {
		return handlers.ErrBadCSRF
	}
	return nil
}
