package util

import (
	"net/http"
)

const (
	JsonContentType = "application/json"
	FormContentType = "application/x-www-form-urlencoded"
)

var WebhookContentTypeList = []string{JsonContentType, FormContentType}

// HttpClient client
var HttpClient = &http.Client{
	//Timeout: 30 * time.Second,
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
