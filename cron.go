package main

import (
	"github.com/lin07ux/go-gin-example/models"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	"github.com/robfig/cron"
	"log"
	"time"
)

func main() {
	log.Println("Starting...")

	setting.Setup()
	models.Setup()

	c := cron.New()
	_ = c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllTag...")
		models.CleanAllTag()
	})
	_ = c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllArticle...")
		models.CleanAllArticle()
	})

	c.Start()

	t := time.NewTimer(time.Second * 10)
	for {
		select {
		case <- t.C:
			t.Reset(time.Second * 10)
		}
	}
}