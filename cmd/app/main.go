package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/asliddinberdiev/i_tv_task/internal/app"
)

func main() {
	app := app.NewCreateApp()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := app.Stop(); err != nil {
		log.Fatal(err)
	}
}
