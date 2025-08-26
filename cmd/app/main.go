package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Eatriceeveryday/data-stream-service/api/http"
	"github.com/Eatriceeveryday/data-stream-service/internal/config"
	"github.com/Eatriceeveryday/data-stream-service/internal/emqx"
	"github.com/Eatriceeveryday/data-stream-service/internal/service"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := emqx.ConnectToClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	emqxService := service.NewEmqxService(client, cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := http.NewHandler(emqxService)

	emqxService.StartPublishing(ctx)

	e := echo.New()
	e.PUT("/sensor", handler.ChangeFrequency)

	go func() {
		if err := e.Start(":8080"); err != nil {
			log.Println("server stopped:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("⏹ Stopping...")
	cancel()
	client.Disconnect(250)
}
