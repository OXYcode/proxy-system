package main

import (
	"net/url"
)

//HTTPReq stores httpRequest fileds
type HTTPReq struct {
	Method     string            `json:"method"`
	Proto      string            `json:"proto"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Host       string            `json:"host"`
	Form       url.Values        `json:"form"`
	Trailer    map[string]string `json:"trailer"`
	RemoteAddr string            `json:"remoteAddr"`
	Target     string            `json:"target"`
}
