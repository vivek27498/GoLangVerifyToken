package Controllers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func searchFromArray(inputArray []string, inputField string) bool {
	for i := 0; i < len(inputArray); i++ {
		if inputArray[i] == inputField {
			return true
		}
	}
	return false
}

func VerifyJwt(req map[string]interface{}) (bool, error) {
	if req["authorization"] == "" || req["authorization"] == nil {
		return false, errors.New("missing authorization header")
	}

	const BEARER_SCHEMA = "Bearer "

	authHeader := req["authorization"]
	if authHeader == "" {
		return false, errors.New("missing authorization header")
	}
	authHeaderValue, ok := authHeader.(string)
	if !ok {
		return false, errors.New("invalid authorization header")
	}

	token := strings.Replace(authHeaderValue, BEARER_SCHEMA, "", 1)

	decodedJwt, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		fmt.Println("Could not parse JWT token:", err)
		return false, errors.New("could not parse jwt token")
	}
	claims, ok := decodedJwt.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Could not parse JWT claims")
		return false, errors.New("could not parse jwt claims")
	}

	// fmt.Println("decodedJwt", claims["scope"].(string))

	scopeArray := strings.Split(claims["scope"].(string), " ")
	fmt.Println("Scope array:", scopeArray)
	fmt.Println("API Scope:", req["servicename"])

	searchRes := searchFromArray(scopeArray, req["servicename"].(string))
	fmt.Println("Search result:", searchRes)
	if !searchRes {
		fmt.Println("invalid scope in jwt token")
		return false, errors.New("invalid scope in jwt token")
	}

	return true, nil
}
