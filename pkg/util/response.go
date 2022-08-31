package util

import (
	"net/http"
)

const (
	JsonContentType = "application/json"
	FormContentType = "application/x-www-form-urlencoded"
)

var WebhookContentTypeList = []string{JsonContentType, FormContentType}

func ValidateContentType(contentType string) bool {
	if contentType == "" {
		return false
	}
	for _, v := range WebhookContentTypeList {
		if v == contentType {
			return true
		}
	}
	return false
}

// HttpClient client
var HttpClient = &http.Client{
	//Timeout: 30 * time.Second,
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
