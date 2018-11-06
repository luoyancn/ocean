package common

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

const (
	CONTENT_TYPE                 = "Content-Type"
	APPLICATION_JSON             = "application/json"
	X_AUTH_TOKEN                 = "X-Auth-Token"
	X_OPENSTACK_NOVA_API_VERSION = "X-Openstack-Nova-Api-Version"
)

var HEADERS = map[string]string{CONTENT_TYPE: APPLICATION_JSON}

func request(method string, url string, headers map[string]string,
	reqdata *string) (*http.Response, error) {
	var reqbuf io.Reader

	if nil != reqdata && "" != *reqdata {
		reqbuf = bytes.NewBufferString(*reqdata)
	} else {
		reqbuf = nil
	}

	req, _ := http.NewRequest(strings.ToUpper(method), url, reqbuf)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	transport := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: transport}
	return client.Do(req)
}

func Get(url string, headers map[string]string) (*http.Response, error) {
	return request("GET", url, headers, nil)
}

func Put(url string, headers map[string]string, reqdata *string) (*http.Response, error) {
	return request("PUT", url, headers, reqdata)
}

func Delete(url string, headers map[string]string) (*http.Response, error) {
	return request("DELETE", url, headers, nil)
}

func Post(url string, headers map[string]string, reqdata *string) (*http.Response, error) {
	return request("POST", url, headers, reqdata)
}
