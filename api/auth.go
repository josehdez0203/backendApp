package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/josehdez0203/backendApp/logger"
)

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

type jwtUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	jwt.RegisteredClaims
}

func (j *Auth) GenerateTokenPair(user *jwtUser) (TokenPairs, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set the Claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = j.Audience
	claims["iss"] = j.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"
	// Set the TokenExpiry
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

	// Create a signed token
	signedAccessToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		logger.L_Error("❌ " + err.Error())
		return TokenPairs{}, err
	}

	// Create refresh token and set Claims
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()
	// set the expirity of refresh Token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()
	// /create  a signed refresh Token
	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create TokenPairs and populate with signed Tokens
	tokenPairs := TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}
	// Rerurn TokenPairs
	return tokenPairs, nil
}

// func (j *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
// 	return &http.Cookie{
// 		Name:     j.CookieName,
// 		Path:     j.CookiePath,
// 		Value:    refreshToken,
// 		Expires:  time.Now().Add(j.RefreshExpiry),
// 		MaxAge:   int(j.RefreshExpiry.Seconds()),
// 		SameSite: http.SameSiteStrictMode,
// 		Domain:   j.CookieDomain,
// 		HttpOnly: true,
// 		Secure:   true,
// 	}
// }
//
// func (j *Auth) GetExpiredRefreshCookie() *http.Cookie {
// 	return &http.Cookie{
// 		Name:     j.CookieName,
// 		Path:     j.CookiePath,
// 		Value:    "",
// 		Expires:  time.Unix(0, 0),
// 		MaxAge:   -1,
// 		SameSite: http.SameSiteStrictMode,
// 		Domain:   j.CookieDomain,
// 		HttpOnly: true,
// 		Secure:   true,
// 	}
// }

func (j *Auth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	w.Header().Add("Vary", "Authorization")
	// Get auth Header
	authHeader := r.Header.Get("Authorization")
	// sanity check
	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}
	// split the header
	headersParts := strings.Split(authHeader, " ")
	if len(headersParts) != 2 {
		return "", nil, errors.New("invalid auth header")
	}
	// check to see if we have the word Bearer
	if headersParts[0] != "Bearer" {
		return "", nil, errors.New("invalid auth header")
	}

	token := headersParts[1]
	// declare empty claims
	claims := &Claims{}
	// parse token
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}
	if claims.Issuer != j.Issuer {
		return "", nil, errors.New("invalid issuer")
	}
	return token, claims, nil
}
