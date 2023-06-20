package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thedevsaddam/gojsonq/v2"
	"github.com/zdz1715/webhook/config"
	"github.com/zdz1715/webhook/global"
	"github.com/zdz1715/webhook/middleware"
	"github.com/zdz1715/webhook/pkg/engine"
	"github.com/zdz1715/webhook/pkg/util"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var jsonParse = gojsonq.New()

type RespInfo struct {
	Status       string      `json:"status"`
	StatusCode   int         `json:"statusCode"`
	ResponseBody interface{} `json:"response_body"`
}

type ReqInfo struct {
	URL    string      `json:"url"`
	Method string      `json:"method"`
	Body   string      `json:"body"`
	Header http.Header `json:"header"`
}

type Req struct {
	URL string `form:"url"`
}

func Handle(c *gin.Context) {
	uuid := strings.ToLower(c.Param("uuid"))
	webhooks := config.Init().Webhooks
	webhook, ok := webhooks[uuid]

	wReq := &Req{}

	_ = c.BindQuery(wReq)

	if len(wReq.URL) > 0 {
		webhook.URL = strings.TrimFunc(wReq.URL, func(r rune) bool {
			return r == '"' || r == '\''
		})
	}

	loggerInfo := global.WebhookLogger.Info().
		Str(middleware.ReqIDKey, middleware.GetReqID(c)).
		Str("webhook_uuid", uuid)

	loggerError := global.WebhookLogger.Error().
		Str(middleware.ReqIDKey, middleware.GetReqID(c)).
		Str("webhook_uuid", uuid)

	name := "webhook/" + uuid

	if !ok {
		msg := name + " is not found"
		loggerError.Msg(msg)
		c.JSON(http.StatusOK, &util.Response{Code: http.StatusNotFound, Message: msg})
		return
	}

	if len(webhook.URL) == 0 {
		msg := name + ".url is empty"
		loggerError.Msg(msg)
		c.JSON(http.StatusOK, &util.Response{Code: http.StatusBadRequest, Message: msg})
		return
	}

	if len(webhook.ContentType) == 0 {
		webhook.ContentType = util.JsonContentType
	}

	if !util.ValidateContentType(webhook.ContentType) {
		msg := fmt.Sprintf("%s.contentType is unsupported: %s, Supported list: %v", name, webhook.ContentType,
			util.WebhookContentTypeList)
		loggerError.Msg(msg)
		c.JSON(http.StatusOK, &util.Response{Code: http.StatusBadRequest, Message: msg})
		return
	}

	reqBody, _ := c.GetRawData()

	vars := bindVars(&webhook, c, reqBody)

	loggerInfo = loggerInfo.Interface("vars", vars)

	strBody, body, err := makeBody(c, &webhook, vars)
	if err != nil {
		loggerError.Str("json str", strBody).Msg(err.Error())
		c.JSON(http.StatusOK, &util.Response{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	req, err := makeRequest(body, c, &webhook, vars)

	if err != nil {
		loggerError.Msg(err.Error())
		c.JSON(http.StatusOK, &util.Response{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	reqInfo := &ReqInfo{
		URL:    req.URL.String(),
		Method: req.Method,
		Body:   strBody,
		Header: req.Header,
	}

	loggerInfo = loggerInfo.Interface("request", reqInfo)

	clientConfig := &config.Client{
		Timeout: config.Init().Client.Timeout,
	}

	if webhook.Client != nil && webhook.Client.Timeout > 0 {
		clientConfig.Timeout = webhook.Client.Timeout
	}
	respInfo, err := sendRequest(req, clientConfig)

	if err != nil {
		loggerError.Msg(err.Error())
		c.JSON(http.StatusOK, &util.Response{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	loggerInfo.Interface("response", respInfo).Msg("forwarding succeeded")
	c.JSON(http.StatusOK, &util.Response{Code: http.StatusOK, Message: "", Data: map[string]interface{}{
		"request":  reqInfo,
		"response": respInfo,
	}})
}

func sendRequest(req *http.Request, client *config.Client) (*RespInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	req.WithContext(ctx)

	defer cancel()

	resp, err := util.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &RespInfo{
		Status:       resp.Status,
		StatusCode:   resp.StatusCode,
		ResponseBody: string(b),
	}, nil
}

func parseVarString(c *gin.Context, text string, vars engine.Vars) string {
	v, err := global.Engine.Render(text, vars)
	if err != nil {
		log.Error().Str(middleware.ReqIDKey, middleware.GetReqID(c)).Msg(err.Error())
		return text
	}
	return v
}

func makeBody(c *gin.Context, webhook *config.Webhook, vars engine.Vars) (string, io.Reader, error) {
	var (
		reqBody io.Reader
		strBody string
	)
	switch webhook.ContentType {
	case util.FormContentType:
		if len(webhook.Body.Form) > 0 {
			postForm := url.Values{}
			for k, v := range webhook.Body.Form {
				postForm.Add(k, parseVarString(c, v, vars))
			}
			strBody = postForm.Encode()
			reqBody = strings.NewReader(strBody)
		}
	case util.JsonContentType:
		if len(webhook.Body.Json) > 0 {
			var bJson interface{}
			err := json.Unmarshal([]byte(webhook.Body.Json), &bJson)
			if err != nil {
				return "", reqBody, fmt.Errorf("config body.json unmarshal failed: %s", err.Error())
			}
			bJsonByte, err := json.Marshal(bJson)
			if err != nil {
				return "", reqBody, fmt.Errorf("config body.json marshal failed: %s", err.Error())
			}
			bJsonStr := string(bJsonByte)
			strBody = parseVarString(c, bJsonStr, vars)
			var jsonObj interface{}
			err = json.Unmarshal([]byte(strBody), &jsonObj)
			if err != nil {
				return strBody, reqBody, fmt.Errorf("body.json after parse unmarshal failed: %s", err.Error())
			}
			reqBody = strings.NewReader(strBody)
		}
	}
	return strBody, reqBody, nil
}

func makeRequest(body io.Reader, c *gin.Context, webhook *config.Webhook, vars engine.Vars) (*http.Request, error) {

	req, err := http.NewRequest(webhook.Method, webhook.URL, body)
	if err != nil {
		return nil, err
	}

	// add header
	for _, headerKV := range webhook.Header {
		req.Header.Add(headerKV.Name, parseVarString(c, headerKV.Value, vars))
	}

	// add query
	if len(webhook.Query) > 0 {
		q := req.URL.Query()
		for _, queryKV := range webhook.Query {
			q.Add(queryKV.Name, parseVarString(c, queryKV.Value, vars))
		}
		req.URL.RawQuery = q.Encode()
	}
	// set and cover header: Content-Type
	req.Header.Set("Content-Type", webhook.ContentType)

	return req, nil
}

func bindVars(webhook *config.Webhook, c *gin.Context, body []byte) engine.Vars {
	data := make(engine.Vars, len(webhook.Vars))
	bodyStr := string(body)
	for name, v := range webhook.Vars {
		nameTrim := strings.TrimSpace(name)
		data[nameTrim] = v.Value

		if len(v.Key) == 0 {
			continue
		}

		key := strings.TrimSpace(v.Key)

		switch v.From {
		case config.WebhookVarFromQuery:
			data[nameTrim] = c.Query(key)
		case config.WebhookVarFromHeader:
			data[nameTrim] = c.GetHeader(key)
		case config.WebhookVarFromBody:
			switch c.ContentType() {
			case util.JsonContentType:
				data[nameTrim] = jsonParse.FromString(bodyStr).Find(key)
			default:
				data[nameTrim] = c.PostForm(key)
			}
		}
	}

	return data
}
