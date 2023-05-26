package Controllers

import (
	"fmt"
	"net/http"
	"strings"
)

func SetVerifyRes(wri http.ResponseWriter, req *http.Request) {

	var authHeader string = req.Header.Get("authorization")
	if req.Header.Get("x-rbl-auth") != "" {
		authHeader = req.Header.Get("x-rbl-auth")
	}

	fmt.Println("req.Body", req.)

	serviceName := req.Header.Get("x-service-name")
	remoteHeader := ""
	xForwardedFor := req.Header.Get(remoteHeader)
	xRealIp := req.Header.Get("x-real-ip")
	ip := req.RemoteAddr
	body := req.Body
	method := req.Method
	headers := req.Header
	url := strings.Split(req.URL.String(), "verifyToken")[1]

	request := map[string]interface{}{
		"authorization": authHeader,
		"servicename":   serviceName,
		"xForwardedFor": xForwardedFor,
		"xRealIp":       xRealIp,
		"ip":            ip,
		"body":          body,
		"method":        method,
		"headers":       headers,
		"url":           url,
	}

	fmt.Println("request URL SPLIT", request["url"])

	finalRes, err := VerifyJwt(request)
	if !finalRes && err != nil && err.Error() == "missing authorization header" {
		panic(err.Error())
	}

	fmt.Println("finalRes", finalRes)

	if finalRes {
		callAxiosResponse, _ := CallAxios(request)
		fmt.Println(callAxiosResponse)
	}

}
