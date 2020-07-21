package main

import (
	"net/http"
	"net/url"
)

//HTTPReq stores httpRequest fileds
type HTTPReq struct {
	Method     string      `json:"method"`
	Proto      string      `json:"proto"`
	Header     http.Header `json:"header"`
	Body       string      `json:"body"`
	Host       string      `json:"host"`
	Form       url.Values  `json:"form"`
	Trailer    http.Header `json:"trailer"`
	RemoteAddr string      `json:"remoteAddr"`
	Target     string      `json:"target"`
}
