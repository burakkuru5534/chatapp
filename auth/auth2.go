package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"net/http"
	"time"
)

type Auth struct {
	*jwtauth.JWTAuth
	devMode bool
}

func NewAuth(secret string, devMode bool) *Auth {
	return &Auth{jwtauth.New("HS256", []byte(secret), nil), devMode}
}

func (t *Auth) NewTokenString(tc *TokenClaims) string {
	mc := t.mapClaims(tc)
	_, tokenString, _ := t.Encode(mc)
	return tokenString
}

func (t *Auth) mapClaims(tc *TokenClaims) jwt.MapClaims {
	var duration time.Duration
	if t.devMode {
		duration = time.Hour * 72
	} else {
		duration = time.Hour * 24
	}

	mc := jwt.MapClaims{
		"exp":          jwtauth.ExpireIn(duration),
		"userID":       tc.UserID,
		"roles":        tc.Roles,
		"userName":     tc.UserName,
		"userEmail":    tc.UserEmail,
		"userFullName": tc.UserFullName,
		"userType":     tc.UserType,
	}

	return mc
}

type TokenClaims struct {
	UserID       int64    `json:"userID"`
	Roles        []string `json:"roles"`
	UserName     string   `json:"userName"`
	UserEmail    string   `json:"userEmail"`
	UserFullName string   `json:"userFullName"`
	UserType     string   `json:"userType"`
}

func TokenClaimsFromRequest(r *http.Request) *TokenClaims {
	_, claims, _ := jwtauth.FromContext(r.Context())
	tc := TokenClaims{
		UserID:       int64(claims["userID"].(float64)),
		Roles:        InterfaceArrayToStringArray(claims["roles"].([]interface{})),
		UserName:     claims["userName"].(string),
		UserEmail:    claims["userEmail"].(string),
		UserFullName: claims["userFullName"].(string),
		UserType:     claims["userType"].(string),
	}

	return &tc
}

func (tc *TokenClaims) IsAdmin() bool {
	return StringInSlice("admin", tc.Roles)
}

func (tc *TokenClaims) HasRole(role string) bool {
	return StringInSlice(role, tc.Roles)
}

func InterfaceArrayToStringArray(t []interface{}) []string {
	s := make([]string, len(t))
	for i, v := range t {
		s[i] = fmt.Sprint(v)
	}

	return s
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
