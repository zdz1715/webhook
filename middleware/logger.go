package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/zdz1715/webhook/global"
	"github.com/zdz1715/webhook/pkg/util"
	"io/ioutil"
	"time"
)

//func WebhookLogPrepare() gin.HandlerFunc {
//	return func(context *gin.Context) {
//		glo
//	}
//}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 提取body 并塞回去
		body, err := c.GetRawData()
		if err == nil {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		if raw != "" {
			path = path + "?" + raw
		}

		// Process request
		c.Next()

		cost := time.Now().Sub(start)

		logger := global.AccessLogger.Info().
			Str(ReqIDKey, GetReqID(c)).
			Str("client_ip", c.ClientIP()).
			Str("time", time.Now().Format(time.RFC3339)).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("proto", c.Request.Proto).
			Int("status_code", c.Writer.Status()).
			Str("latency", cost.String()).
			Str("user_agent", c.Request.UserAgent()).
			Str("query", raw).
			Interface("header", c.Request.Header)

		if c.ContentType() == util.JsonContentType {
			logger = logger.RawJSON("body", body)
		} else {
			logger = logger.Interface("body", string(body))
		}

		logger.Str("error", c.Errors.ByType(gin.ErrorTypePrivate).String()).Msg("")

	}
}
