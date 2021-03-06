// Package rest allows for quick and easy access any REST or REST-like API.
package rest

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// Method contains the supported HTTP verbs.
type Method string

const (
	Get    Method = "GET"
	Post   Method = "POST"
	Put    Method = "PUT"
	Patch  Method = "PATCH"
	Delete Method = "DELETE"
)

// Request holds the request to an API Call.
type Request struct {
	Method         Method
	BaseURL        string // e.g. https://api.sendgrid.com
	RequestHeaders map[string]string
	QueryParams    map[string]string
	RequestBody    []byte
}

// Response holds the response from an API call.
type Response struct {
	StatusCode      int                 // e.g. 200
	ResponseBody    string              // e.g. {"result: success"}
	ResponseHeaders map[string][]string // e.g. map[X-Ratelimit-Limit:[600]]
}


// BuildRequestObject creates the HTTP request object.
func BuildRequestObject(request Request) (*http.Request, error) {
	req, err := http.NewRequest(string(request.Method), request.BaseURL, bytes.NewBuffer(request.RequestBody))

	// adds query paramaters to the URL.
	if len(request.QueryParams) != 0 {
		q := req.URL.Query()
		for key, value := range request.QueryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	for key, value := range request.RequestHeaders {
		req.Header.Set(key, value)
	}

	return req, err
}

// MakeRequest makes the API call.
func MakeRequest(req *http.Request) (*http.Response, error) {
	var Client = &http.Client{
		Transport: http.DefaultTransport,
	}
	res, err := Client.Do(req)
	return res, err
}

// BuildResponse builds the response struct.
func BuildResponse(res *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	response := Response{
		StatusCode:      res.StatusCode,
		ResponseBody:    string(body),
		ResponseHeaders: res.Header,
	}
	return &response, nil
}

// API is the main interface to the API.
func API(request Request) (*Response, error) {

	// Build the HTTP request object.
	req, err := BuildRequestObject(request)
	if err != nil {
		return nil, err
	}

	// Build the HTTP client and make the request.
	res, err := MakeRequest(req)
	if err != nil {
		return nil, err
	}

	// Build Response object.
	response, err := BuildResponse(res)
	if err != nil {
		return nil, err
	}

	return response, nil
}
