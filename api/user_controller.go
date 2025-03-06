package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/josehdez0203/backendApp/logger"

	"github.com/golang-jwt/jwt/v4"
	"github.com/josehdez0203/backendApp/models"
	"golang.org/x/crypto/bcrypt"
)

func (app application) newUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := app.readJSON(w, r, &user)
	if err != nil {
		logger.L_Error("Error :" + err.Error())
		app.errorJSON(w, err, http.StatusBadRequest)
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		logger.L_Error("Error " + err.Error())
		app.errorJSON(w, err, http.StatusBadRequest)
	}
	user.Password = string(hashedPass)
	newUser, err := app.DB.AddUser(user)
	if err != nil {
		logger.L_Error("Error " + err.Error())
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}
	logger.L_Info("New User created 👍")
	app.writeJSON(w, http.StatusOK, newUser)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload
	// user := models.User{}

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	// validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	// user, err := Get
	if err != nil {
		app.errorJSON(w, errors.New("Invalid credenctials"), http.StatusBadRequest)
		return
	}
	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("Invalid credenctials"), http.StatusBadRequest)
		return
	}
	// create a jwt user

	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	logger.L_Info(tokens.Token)
	// refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	// http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

//	func (app application) refreshToken(w http.ResponseWriter, r *http.Request) {
//		logger.L_Info("refreshToken")
//		if len(r.Cookies()) == 0 {
//			logger.L_Error("Sin credenciales")
//			app.errorJSON(w, errors.New("Sin credenciales"), http.StatusUnauthorized)
//			return
//		}
//		for _, cookie := range r.Cookies() {
//			logger.L_Info(cookie.Name)
//			if cookie.Name == app.auth.CookieName {
//				claims := &Claims{}
//				refreshToken := cookie.Value
//
//				// parse the token to get the claims
//				_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
//					return []byte(app.JWTSecret), nil
//				})
//				if err != nil {
//					logger.L_Error("No authorized ❌")
//					app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
//					return
//				}
//				// get user id from de token claims
//				userID, err := strconv.Atoi(claims.Subject)
//				if err != nil {
//					app.errorJSON(w, errors.New("unknow user"), http.StatusUnauthorized)
//					return
//				}
//
//				user, err := app.DB.GetUserById(userID)
//				if err != nil {
//					app.errorJSON(w, errors.New("unknow user"), http.StatusUnauthorized)
//					return
//				}
//
//				u := jwtUser{
//					ID:        user.ID,
//					FirstName: user.FirstName,
//					LastName:  user.LastName,
//				}
//
//				tokenPairs, err := app.auth.GenerateTokenPair(&u)
//				if err != nil {
//					app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
//					return
//				}
//				http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))
//
//				app.writeJSON(w, http.StatusOK, tokenPairs)
//			}
//		}
//	}
func (app application) refreshToken(w http.ResponseWriter, r *http.Request) {
	logger.L_Info("refreshToken")
	var requestPayload struct {
		Token string `json:"token"`
	}
	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &requestPayload)
	if err != nil {
		logger.L_Error(err.Error())
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	logger.L_Info(requestPayload.Token)
	claims := &Claims{}
	refreshToken := requestPayload.Token

	// parse the token to get the claims
	_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.JWTSecret), nil
	})
	if err != nil {
		logger.L_Error("No authorized ❌")
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}
	// get user id from de token claims
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		logger.L_Error("No authorized ❌" + err.Error())
		app.errorJSON(w, errors.New("unknow user"), http.StatusUnauthorized)
		return
	}

	user, err := app.DB.GetUserById(userID)
	if err != nil {
		logger.L_Error("No authorized ❌" + err.Error())
		app.errorJSON(w, errors.New("unknow user"), http.StatusUnauthorized)
		return
	}

	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	tokenPairs, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
		return
	}
	// http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))

	app.writeJSON(w, http.StatusOK, tokenPairs)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	logger.L_Info("logging out")
	// http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}
