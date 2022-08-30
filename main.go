package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/zdz1715/webhook/config"
	"github.com/zdz1715/webhook/pkg/util"
	"github.com/zdz1715/webhook/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	pflag.Parse()

	// 载入config
	cfg := config.Init()

	// 开启http服务
	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Application.Host, cfg.Application.Port),
		Handler:        router.Init(),
		ReadTimeout:    cfg.Application.ReadTimeout,
		WriteTimeout:   cfg.Application.WriteTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	log.Info().Msgf("root: %s, listen: %s:%d", util.GetExecRootDir(), cfg.Application.Host, cfg.Application.Port)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msg("Server exiting")

}
