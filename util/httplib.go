package util

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	GET  = "GET"
	POST = "POST"
)

//Response.Status definition
const (
	WRONG = iota - 2 //-2, request wrong
	FAIL             //-1, request fail
)

//POST TYPE
const (
	FORM = iota + 1 //x-www-form-urlencoded
	DATA            //multipart/form-data
	JSON            //application/json
	XML             //text/xml
)

type Hashmap map[string]string

type Request struct {
	Method           string //support GET & POST only for now
	Url              string
	Headers          Hashmap
	Ssl              bool        //whether https request
	ConnectTimeout   int         //connection timeout in millisecond, default 200ms
	ReadWriteTimeout int         //read write timeout in millisecond, default 200ms
	Type             int         //post type
	Params           interface{} //different post params according to post type
}

type Response struct {
	Status        int
	Message, Body string
}

/*Call method param list*/
const (
	SSL = iota
	POST_TYPE
	PARAM
	HEADER
	CONNECTTIMEOUT
	READWRITETIMEOUT
)

/**
* desc singel request
* @param method request type
* @param url
* @param ssl	whether https request.
* @param type	post type
* @param params  only valid in POST method
* @param headers such as {"cookie":"a=b; c=d", "user-agent":"go", "x-id":123}
* @param ConnectTimeout
* @param ReadWriteTimeout
* @return *Response
**/
func Call(req Request) (res *Response) {
	res = &Response{}
	if req, msg, ok := check(req); ok {
		call(req, res, nil, nil)
	} else {
		res.Status = WRONG
		res.Message = msg
	}
	return
}

var max_concurrency = 1000   //multi-request max concurrency
var default_concurrency = 10 //default concurrency

/**
* desc multi-request
* @param reqs request definitions
* @param concurrency
* @return map[interface{}]*Response
**/
func MultiCall(reqs map[interface{}]Request, arg ...int) map[interface{}]*Response {
	res := make(map[interface{}]*Response)
	if len(reqs) == 0 {
		return res
	}
	concurrency := default_concurrency
	if arg != nil {
		num := arg[0]
		switch {
		case num > max_concurrency:
			num = max_concurrency
			fallthrough
		default:
			concurrency = num
		}
	}
	ch := make(chan int, concurrency)
	wg := &sync.WaitGroup{}
	for idx, req := range reqs {
		if req, msg, ok := check(req); ok {
			res[idx] = &Response{}
			ch <- 1
			wg.Add(1)
			go call(req, res[idx], ch, wg)
		} else {
			res[idx] = &Response{Status: WRONG, Message: msg}
		}
	}
	wg.Wait()
	close(ch)
	return res
}

/**
* desc check request params
* @return Request
**/
func check(req Request) (*Request, string, bool) {
	request := &Request{
		Headers: req.Headers,
		Url:     req.Url,
	}
	if req.Method != GET && req.Method != POST {
		return nil, "wrong method", false
	}
	request.Method = req.Method

	if req.Method == POST {
		switch req.Type {
		case FORM:
		case DATA:
		case JSON:
		case XML:
		default:
			return nil, "post type invalid", false
		}
		request.Type = req.Type
		request.Params = req.Params
	}

	if strings.Contains(req.Url, "https://") {
		request.Ssl = true
	} else {
		request.Ssl = req.Ssl
	}
	if req.ConnectTimeout <= 0 {
		request.ConnectTimeout = 200
	} else {
		request.ConnectTimeout = req.ConnectTimeout
	}

	if req.ReadWriteTimeout <= 0 {
		request.ReadWriteTimeout = 200
	} else {
		request.ReadWriteTimeout = req.ReadWriteTimeout
	}
	return request, "", true
}

/**
* desc call go standard net tool
**/
func call(req *Request, res *Response, ch chan int, wg *sync.WaitGroup) {
	defer func() {
		if ch != nil {
			<-ch
		}
		if wg != nil {
			wg.Done()
		}
	}()
	request, err := http.NewRequest(req.Method, req.Url, nil) //strings.NewReader(params.Encode()))
	if err != nil {
		res.Status = FAIL
		res.Message = err.Error()
		return
	}
	for idx, val := range req.Headers {
		request.Header.Set(idx, val)
	}
	if req.Method == POST {
		switch req.Type {
		case FORM:
			params, _ := req.Params.(Hashmap)
			param := url.Values{}
			for idx, val := range params {
				param.Add(idx, val)
			}
			body := param.Encode()
			request.ContentLength = int64(len(body))
			request.Body = ioutil.NopCloser(strings.NewReader(body))
			request.Header.Set("content-type", "application/x-www-form-urlencoded")
		case XML:
			param, _ := req.Params.(string)
			request.ContentLength = int64(len(param))
			request.Header.Set("content-type", "text/xml")
			request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(param)))
		case JSON:
			param, _ := req.Params.(string)
			request.ContentLength = int64(len(param))
			request.Header.Set("content-type", "application/json")
			request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(param)))
		}
	}
	tr := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(network, addr, time.Millisecond*time.Duration(req.ConnectTimeout))
			if err != nil {
				return nil, err
			}
			conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(req.ReadWriteTimeout)))
			return conn, nil
		},
	}
	if req.Ssl {
		tr.TLSHandshakeTimeout = time.Millisecond * time.Duration(req.ConnectTimeout)
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{
		Transport: tr,
	}
	response, err := client.Do(request)
	if err != nil {
		res.Status = FAIL
		res.Message = err.Error()
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		res.Status = FAIL
		res.Message = err.Error()
		return
	}
	res.Status = response.StatusCode
	res.Body = string(body)
	return
}
