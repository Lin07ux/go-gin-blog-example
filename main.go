package main

import (
	"fmt"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	"github.com/lin07ux/go-gin-example/routers"
	"net/http"
)

func main() {
	s := &http.Server{
		Addr:              fmt.Sprintf(":%d", setting.HttpPort),
		Handler:           routers.InitRouter(),
		ReadTimeout:       setting.ReadTimeout,
		WriteTimeout:      setting.WriteTimeout,
		MaxHeaderBytes:    1 << 20,
	}

	_ = s.ListenAndServe()
}
