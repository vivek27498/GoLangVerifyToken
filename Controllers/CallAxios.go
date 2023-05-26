package Controllers

import (
    "fmt"
    "net/http"
    "io/ioutil"
	"bytes"
)

func CallAxios(input map[string]interface{}) ([]byte, error) {

    fmt.Println("input",input)

    url := fmt.Sprintf("%s%s", "https://cloudapiuat.ratnakarbank.in", input["url"].(string))
    fmt.Println("method",input["method"].(string))
    fmt.Println("url",url)

	body, ok := input["body"].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid method: %v", input["body"])
	}

    fmt.Println("body",body)

    req, err := http.NewRequest(input["method"].(string), url, bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }

    for key, value := range input["headers"].(map[string]string) {
        req.Header.Set(key, value)
    }
    req.Header.Del("Content-Length")

    fmt.Println("req",req)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    fmt.Println("responseBody",responseBody)

    return body, nil
}