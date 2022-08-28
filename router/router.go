package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/thedevsaddam/gojsonq/v2"
	"github.com/zdz1715/webhook/config"
	"github.com/zdz1715/webhook/library/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"text/template"
	"time"
)

func Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	//gin.DisableConsoleColor()
	g := gin.Default()
	initRouter(g)
	return g
}

func initRouter(g *gin.Engine) {
	logCfg := config.Init().Log
	wlog = log.NewRotate(logCfg.Level, path.Join(logCfg.DIR, "webhooks.log"),
		logCfg.Rotate.MaxSize, logCfg.Rotate.MaxBackups)

	g.Any("/webhook/:uuid", handleWebhook)

}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var wlog zerolog.Logger
var httpClient = &http.Client{
	//Timeout:
}

func bindWebhookVars(vars map[string]config.WebhookVar, c *gin.Context, rawData []byte) map[string]interface{} {
	data := make(map[string]interface{}, len(vars))
	for name, v := range vars {
		data[name] = v.Value

		if len(v.Key) == 0 {
			continue
		}

		switch v.From {
		case config.WebhookVarFromQuery:
			data[name] = c.Query(v.Key)
		case config.WebhookVarFromHeader:
			data[name] = c.GetHeader(v.Key)
		case config.WebhookVarFromBody:
			switch c.ContentType() {
			case "application/json":
				data[name] = gojsonq.New().FromString(string(rawData)).Find(v.Key)
			default:
				data[name] = c.PostForm(v.Key)
			}
		}
	}

	return data
}

func parseVarString(text string, data map[string]interface{}) string {
	t, err := template.New("").Parse(text)
	if err != nil {
		return text
	}
	var buf bytes.Buffer

	err = t.Execute(&buf, data)
	if err != nil {
		return text
	}
	return buf.String()
}

func handleWebhook(c *gin.Context) {
	uuid := strings.ToLower(c.Param("uuid"))
	webhooks := config.Init().Webhooks
	webhook, ok := webhooks[uuid]

	d, _ := c.GetRawData()

	// init log
	wlog = wlog.With().
		Str("client_ip", c.ClientIP()).
		Str("url", c.Request.URL.String()).
		Str("user_agent", c.GetHeader("User-Agent")).
		Str("method", c.Request.Method).
		Str("body", string(d)).
		Interface("header", c.Request.Header).
		Logger()

	wlog.Debug().Msg(fmt.Sprintf("uuid: %s", uuid))
	name := "webhook/" + uuid
	if !ok {
		msg := name + " is not found"
		wlog.Error().Msg(msg)
		c.JSON(http.StatusOK, &Response{Code: http.StatusNotFound, Message: msg})
		return
	}

	if len(webhook.URL) == 0 {
		msg := name + ".url is empty"
		wlog.Error().Msg(msg)
		c.JSON(http.StatusOK, &Response{Code: http.StatusBadRequest, Message: msg})
		return
	}
	if len(webhook.ContentType) == 0 {
		webhook.ContentType = "application/json"
	}

	webhookVars := bindWebhookVars(webhook.Vars, c, d)

	wlog.Debug().Msg(fmt.Sprintf("webhook vars: %+v", webhookVars))

	var reqBody io.Reader
	// add body
	if len(webhook.Body) > 0 {
		switch webhook.ContentType {
		case "application/json":
			jsonBody, err := json.Marshal(webhook.Body)
			if err != nil {
				wlog.Error().Err(err)
				c.JSON(http.StatusOK, &Response{Code: http.StatusBadRequest, Message: err.Error()})
				return
			}
			strBody := parseVarString(string(jsonBody), webhookVars)
			reqBody = strings.NewReader(strBody)
		case "application/x-www-form-urlencoded":
			postForm := url.Values{}
			for k, v := range webhook.Body {
				if vStr, ok := v.(string); ok {
					postForm.Add(k, parseVarString(vStr, webhookVars))
				}
			}
			reqBody = strings.NewReader(postForm.Encode())
		default:
			msg := name + ".contentType is unsupported: " + webhook.ContentType
			wlog.Error().Msg(msg)
			c.JSON(http.StatusOK, &Response{Code: http.StatusBadRequest, Message: msg})
			return
		}
	}

	req, err := http.NewRequest(webhook.Method, webhook.URL, reqBody)
	if err != nil {
		wlog.Error().Err(err)
		c.JSON(http.StatusOK, &Response{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	// add header
	for _, headerKV := range webhook.Header {
		req.Header.Add(headerKV.Name, parseVarString(headerKV.Value, webhookVars))
	}
	// add query
	if len(webhook.Query) > 0 {
		q := req.URL.Query()
		for _, queryKV := range webhook.Query {
			q.Add(queryKV.Name, parseVarString(queryKV.Value, webhookVars))
		}
		req.URL.RawQuery = q.Encode()
	}

	// set and cover header: Content-Type
	req.Header.Set("Content-Type", webhook.ContentType)

	timeout := config.Init().Client.Timeout
	retryCount := config.Init().Client.RetryCount
	retrySleepTime := config.Init().Client.RetrySleepTime

	if webhook.Client != nil && webhook.Client.Timeout > 0 {
		timeout = webhook.Client.Timeout
	}

	if webhook.Client != nil && webhook.Client.RetryCount >= 0 {
		retryCount = webhook.Client.RetryCount
	}

	if webhook.Client != nil && webhook.Client.RetrySleepTime > 0 {
		retrySleepTime = webhook.Client.RetrySleepTime
	}

	// httpClient.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	req.WithContext(ctx)
	defer cancel()

	var (
		responseErr error
		resp        *http.Response
	)

	for i := 0; i <= retryCount; i++ {
		resp, responseErr = httpClient.Do(req)
		if responseErr != nil {
			wlog.Error().Err(responseErr).Msg(fmt.Sprintf("retry: %d", i))
			time.Sleep(retrySleepTime)
		} else {
			break
		}
	}

	if responseErr != nil {
		wlog.Error().Err(responseErr)
		c.JSON(http.StatusOK, &Response{Code: http.StatusBadRequest, Message: responseErr.Error()})
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	b, _ := ioutil.ReadAll(resp.Body)
	res := map[string]interface{}{
		"status":        resp.Status,
		"status_code":   resp.StatusCode,
		"response_body": string(b),
	}
	c.JSON(http.StatusOK, &Response{Code: http.StatusOK, Message: "", Data: res})

	wlog.Info().Interface("response", res)
}
