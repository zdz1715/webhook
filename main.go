package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/zdz1715/webhook/config"
	"github.com/zdz1715/webhook/router"
	"net/http"
)

func main() {
	pflag.Parse()

	// 载入config
	cfg := config.Init()

	// 开启http服务
	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Application.Host, cfg.Application.Port),
		Handler:        router.Init(),
		ReadTimeout:    cfg.Application.ReadTimeout,
		WriteTimeout:   cfg.Application.WriteTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}
	fmt.Printf("[HTTP] listen %s:%d \n", cfg.Application.Host, cfg.Application.Port)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
