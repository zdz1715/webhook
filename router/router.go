package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zdz1715/webhook/api/webhook"
	"github.com/zdz1715/webhook/config"
	"github.com/zdz1715/webhook/global"
	"github.com/zdz1715/webhook/middleware"
	"github.com/zdz1715/webhook/pkg/util"
	"io"
	"net/http"
	"os"
	"path"
)

func Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	cfg := config.Init()

	accessLogPath := path.Join(cfg.Log.DIR, "access.log")
	errorLogPath := path.Join(cfg.Log.DIR, "error.log")
	webhookLogPath := path.Join(cfg.Log.DIR, "webhook_reply.log")

	global.AccessLogger = util.NewMultiLevelWriter(
		util.NewLogRotate(accessLogPath, cfg.Log.Rotate.MaxSize, cfg.Log.Rotate.MaxBackups),
		os.Stdout,
	)
	global.WebhookLogger = util.NewLogRotate(webhookLogPath, cfg.Log.Rotate.MaxSize, cfg.Log.Rotate.MaxBackups)

	errorLog := util.NewRotate(errorLogPath, cfg.Log.Rotate.MaxSize, cfg.Log.Rotate.MaxBackups)

	gin.DefaultErrorWriter = io.MultiWriter(errorLog, os.Stderr)

	r := gin.New()

	// 生成unique req id
	r.Use(middleware.ReqID())

	r.Use(middleware.AccessLog())

	r.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		c.JSON(http.StatusInternalServerError, &util.Response{Code: http.StatusInternalServerError, Message: "server error"})
		c.Abort()
	}))

	// registry webhook api
	r.Any("/webhook/:uuid", webhook.Handle)

	return r
}
