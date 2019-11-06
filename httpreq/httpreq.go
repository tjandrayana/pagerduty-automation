package httpreq

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Request struct {
	Url          string
	Path         string
	Method       string
	Headers      map[string]string
	Param        url.Values
	Json         interface{}
	IsJson       bool
	Timeout      time.Duration
	ResultStatus int
}

func NewHTTPRequest() *Request {
	return &Request{
		Headers: map[string]string{},
	}
}

func (r *Request) JSONPost(u *url.URL) (*http.Request, error) {
	u.RawQuery = r.Param.Encode()
	u.Path += r.Path
	link := u.String()

	var err error
	var body []byte

	if reflect.TypeOf(r.Json).Kind() == reflect.String {
		str := r.Json.(string)
		body = []byte(str)
	} else {
		body, err = json.Marshal(r.Json)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(r.Method, link, bytes.NewReader(body))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	r.Headers["Content-Type"] = "application/json"

	fmt.Printf("\nRequest Body: %s\n\n", string(body))

	return req, nil
}

func (r *Request) Post(u *url.URL) (*http.Request, error) {
	if r.IsJson {
		return r.JSONPost(u)
	}
	u.Path += r.Path
	link := u.String()
	form := strings.NewReader(r.Param.Encode())
	req, err := http.NewRequest(r.Method, link, form)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return req, nil
}

func (r *Request) Get(u *url.URL) (*http.Request, error) {
	u.RawQuery = r.Param.Encode()
	u.Path += r.Path
	link := u.String()
	req, err := http.NewRequest(r.Method, link, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return req, nil
}

// DoReq is Last Point to call api
func (r *Request) DoReq() (*[]byte, error) {
	u, err := url.Parse(r.Url)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var req *http.Request
	switch r.Method {
	case "GET", "DELETE":
		req, err = r.Get(u)
	case "POST":
		req, err = r.Post(u)
	case "PUT":
		req, err = r.Put(u)
	}

	if req == nil {
		return nil, fmt.Errorf("Failed create new request")
	}
	if err != nil {
		return nil, err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	// For intermittent EOF error
	req.Close = true

	// set timeout
	// if r.Timeout is defined, use r.Timeout as HTTPReq timeout
	// else take from default
	var timeout time.Duration
	if r.Timeout > 0 {
		timeout = r.Timeout
	} else {
		timeout = 60 * time.Second
	}

	hc := &http.Client{
		Timeout: timeout,
	}

	resp, err := hc.Do(req)
	if err != nil {
		log.Println(resp, err)
		return nil, err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	r.ResultStatus = resp.StatusCode
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		dump, err := httputil.DumpResponse(resp, true)
		log.Println(string(dump), err)
		return &contents, fmt.Errorf("Status Code = %d", resp.StatusCode)
	}

	return &contents, nil
}

func (r *Request) Debug() {
	var content string

	u, err := url.Parse(r.Url)
	if err != nil {
		log.Println("URL is invalid.")
	}

	switch r.Method {
	case http.MethodGet:
		u.RawQuery = r.Param.Encode()
	case http.MethodPost:
		if r.IsJson {
			var body []byte
			var err error
			if reflect.TypeOf(r.Json).Kind() == reflect.String {
				str := r.Json.(string)
				body = []byte(str)
			} else {
				body, err = json.Marshal(r.Json)
				if err != nil {
					log.Println("Json content post is invalid.")
				}
			}

			content = string(body)
		} else {
			content = r.Param.Encode()
		}
	default:
		log.Println("Method is invalid.")
	}

	u.Path += r.Path
	link := u.String()
	headers := r.Headers

	log.Println(link, content, headers)
}

func (r *Request) Put(u *url.URL) (*http.Request, error) {
	if r.IsJson {
		return r.JSONPost(u)
	}
	u.Path += r.Path
	link := u.String()
	form := strings.NewReader(r.Param.Encode())
	req, err := http.NewRequest(r.Method, link, form)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return req, nil
}
