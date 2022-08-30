package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

const (
	ReqIDKey = "req_id"
)

func GetReqID(c *gin.Context) string {
	v, _ := c.Get(ReqIDKey)
	if str, ok := v.(xid.ID); ok {
		return str.String()
	}
	return ""
}

func ReqID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ReqIDKey, xid.New())

		c.Next()

	}
}
