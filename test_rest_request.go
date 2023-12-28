package test_rest_request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

)

var Handler http.Handler



func init(){
	if Handler == nil{
		Handler = http.DefaultServeMux
	}
}

func GetEndpoint(path string, headers map[string]string, responseStruct interface{}) error {

	req := httptest.NewRequest(http.MethodGet, path, nil)

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	w := httptest.NewRecorder()
	Handler.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(res.Status)
	fmt.Print(string(data))

	if res.StatusCode != 200 {
		return fmt.Errorf("%s %s %s", path, res.Status, string(data))
	}
	if responseStruct != nil {
		err = json.Unmarshal(data, responseStruct)
	}

	return err
}

func PostEndpoint(path string, headers map[string]string,  body interface{}, responseStruct interface{}) error {
	return postOrPutEndpoint(http.MethodPost, headers, path, body, responseStruct)
}
func PutEndpoint(path string,  headers map[string]string, body interface{}, responseStruct interface{}) error {
	return postOrPutEndpoint(http.MethodPut, headers, path, body, responseStruct)
}
func DeleteEndpoint(path string,  headers map[string]string, body interface{}, responseStruct interface{}) error {
	return postOrPutEndpoint(http.MethodDelete, headers, path, body, responseStruct)
}
func postOrPutEndpoint(method string, headers map[string]string, path string, body interface{}, responseStruct interface{}) error {

	var req *http.Request

	if values, ok := body.(url.Values); ok{
		var err error
		req, err = http.NewRequest(method, path, strings.NewReader(values.Encode()))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	}else{
		var buf bytes.Buffer
		if body != nil {
			err := json.NewEncoder(&buf).Encode(body)
			if err != nil {
				return nil
			}
		}
		reqBody := buf.String()
		print(reqBody)

		req = httptest.NewRequest(method, path, &buf)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	w := httptest.NewRecorder()
	Handler.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	//fmt.Print(reqBody)
	fmt.Println(res.Status)
	fmt.Print(string(data))

	if res.StatusCode != 200 {
		return fmt.Errorf("%s %s", res.Status, string(data))
	}

	if responseStruct != nil {
		err = json.Unmarshal(data, responseStruct)
	}

	return err
}
