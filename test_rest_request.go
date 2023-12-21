package test_rest_request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func GetEndpoint(path string, responseStruct interface{}, idToken string) error {

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

func PostEndpoint(path string, body interface{}, responseStruct interface{}, idToken string) error {
	return postOrPutEndpoint(http.MethodPost, path, body, responseStruct, idToken)
}
func PutEndpoint(path string, body interface{}, responseStruct interface{}, idToken string) error {
	return postOrPutEndpoint(http.MethodPut, path, body, responseStruct, idToken)
}
func DeleteEndpoint(path string, body interface{}, responseStruct interface{}, idToken string) error {
	return postOrPutEndpoint(http.MethodDelete, path, body, responseStruct, idToken)
}
func postOrPutEndpoint(method string, path string, body interface{}, responseStruct interface{}, idToken string) error {

	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return nil
		}
	}
	reqBody := buf.String()

	req := httptest.NewRequest(method, path, &buf)

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

	fmt.Print(reqBody)
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
