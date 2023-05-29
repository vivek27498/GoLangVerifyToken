package Controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type Config struct {
	Region          string
	CognitoPoolId   string
	CognitoClientId string
}

func New(config *Config) *Config {
	return config
}

func (config *Config) Validate(jwtToken string) error {

	pKey, err := getPublicKeys(config.Region, config.CognitoPoolId)

	if err != nil {
		log.Fatal("Error trying to get Cognito public keys, check your config")
	}

	keySet, _ := jwk.Parse(pKey)

	fmt.Println("keySet", keySet)

	parsedToken, err := jwt.Parse([]byte(jwtToken), jwt.WithKeySet(keySet))

	jsonBytes, err := json.Marshal(parsedToken)
	if err != nil {
		log.Fatal("Error converting keySet to JSON:", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("data", int(data["exp"].(float64)*1000000))

	fmt.Println(parsedToken)

	authTime, _ := parsedToken.Get("exp")
	fmt.Println("authTime", authTime)

	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Failed to load IST time zone:", err)
	}

	istTime := authTime.(time.Time).In(ist)

	fmt.Println("authTime", istTime)

	timestamp := time.Unix(int64(authTime.(float64)), 0)

	fmt.Println("parsedToken", timestamp)

	if err != nil {
		return errors.New("INVALID TOKEN")
	}

	clientId, _ := parsedToken.Get("client_id")
	token_use, _ := parsedToken.Get("token_use")

	if clientId != config.CognitoClientId {
		return errors.New("TOKEN IS FROM A DIFFERENT client_id")
	}

	if parsedToken.Issuer() != fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", config.Region, config.CognitoPoolId) {
		return errors.New("TOKEN IS FROM A DIFFERENT pool_id")
	}

	if token_use != "id" && token_use != "access" {
		return errors.New("TOKEN IS FROM A DIFFERENT source")
	}

	// if time.Now().After(timestamp) {
	// 	return errors.New("TOKEN EXPIRED")
	// }

	return nil
}

func getPublicKeys(region string, cognitoPoolId string) ([]byte, error) {

	var url = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, cognitoPoolId)

	resp, err :=
		http.Get(url)

	if err != nil {
		fmt.Println("Error fetching public keys")
		return nil, errors.New("Error")
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	return body, nil
}
