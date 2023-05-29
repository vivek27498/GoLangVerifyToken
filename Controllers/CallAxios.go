package Controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CallAxios(input map[string]interface{}, body []byte) ([]byte, error) {

	// fmt.Println("input", input)

	url := fmt.Sprintf("%s%s", "http://localhost:9090", input["url"].(string))
	fmt.Println("method", input["method"].(string))
	fmt.Println("url", url)

	req, err := http.NewRequest(input["method"].(string), url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	//convert headers into Map Object
	headers, ok := input["headers"].(http.Header)
	if ok {
		for key, values := range headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	req.Header.Del("Content-Length")

	// fmt.Println("req", req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println("till here")

	fmt.Println("resp", resp)

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
