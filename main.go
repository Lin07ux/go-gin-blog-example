package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	"github.com/lin07ux/go-gin-example/routers"
	"log"
	"syscall"
)

func main() {
	endless.DefaultReadTimeOut = setting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endpoint := fmt.Sprintf(":%d", setting.HttpPort)

	server := endless.NewServer(endpoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
