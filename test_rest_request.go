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

type Client struct{
	Headers map[string]string
}

func init(){
	if Handler == nil{
		Handler = http.DefaultServeMux
	}
}

func (m Client) GetEndpoint(path string, responseStruct interface{}, idToken string) error {

	req := httptest.NewRequest(http.MethodGet, path, nil)

	if idToken != "" {
		req.Header.Set("idToken", idToken)
	}

	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

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

func (m Client) PostEndpoint(path string, body interface{}, responseStruct interface{}, idToken string) error {
	return postOrPutEndpoint(http.MethodPost, m.Headers, path, body, responseStruct, idToken)
}
func (m Client) PutEndpoint(path string, body interface{}, responseStruct interface{}, idToken string) error {
	return postOrPutEndpoint(http.MethodPut, m.Headers, path, body, responseStruct, idToken)
}
func (m Client) DeleteEndpoint(path string, body interface{}, responseStruct interface{}, idToken string) error {
	return postOrPutEndpoint(http.MethodDelete, m.Headers, path, body, responseStruct, idToken)
}
func postOrPutEndpoint(method string, headers map[string]string, path string, body interface{}, responseStruct interface{}, idToken string) error {

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
	if idToken != "" {
		req.Header.Set("idToken", idToken)
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
