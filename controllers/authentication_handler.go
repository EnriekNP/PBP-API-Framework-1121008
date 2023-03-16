package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("Bebasapasaja123!")
var tokenName = "token"

type Claims struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserType int    `json:"user_type"`
	jwt.StandardClaims
}

func generateToken(w http.ResponseWriter, id int, name string, userType int) {
	tokenExpiryTime := time.Now().Add(5 * time.Minute)

	//create claims with user data
	claims := &Claims{
		ID:       id,
		Name:     name,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}

	//encrypt claim to jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//jwtKey :=
	jwtToken, err := token.SignedString(jwtKey)
	if err != nil {
		return
	}

	//set token to cookies
	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    jwtToken,
		Expires:  tokenExpiryTime,
		Secure:   false,
		HttpOnly: true,
	})

}

func resetUserToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   false,
		HttpOnly: true,
	})
}
// func martiniAuthenticate(accessType int) martini.Handler {
// 	return func(c martini.Context) {
// 		w := 
// 		r := c.Request()
// 		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			c.Next()
// 		})

// 		// call the Authenticate function with the correct arguments
// 		middleware := Authenticate(next, accessType)

// 		// call the middleware with the current request and response
// 		middleware.ServeHTTP(w, r)
// 	}
// }
func Authenticate(next http.HandlerFunc, accessType int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isValidToken := validateUserToken(r, accessType)
		if !isValidToken {
			sendUnAuthorizedResponse(w)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func validateUserToken(r *http.Request, accessType int) bool {
	isAccessTokenValid, id, name, userType := validateTokenFromCookies(r)
	fmt.Print(id, name, userType, accessType, isAccessTokenValid)

	if isAccessTokenValid {
		isUserValid := userType == accessType
		if isUserValid {
			return true
		}
	}
	return false
}

func validateTokenFromCookies(r *http.Request) (bool, int, string, int) {
	if cookie, err := r.Cookie(tokenName); err == nil {
		jwtToken := cookie.Value
		accessClaims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(jwtToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return true, accessClaims.ID, accessClaims.Name, accessClaims.UserType
		}
	}
	return false, -1, "", -1
}
