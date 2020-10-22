package main

import (
	"context"
	"fmt"
	"github.com/Logiase/gomirai/bot"
	"github.com/jinzhu/configor"
	"github.com/labstack/echo/v4"
	validator "gopkg.in/go-playground/validator.v9"
	"io"
	"log"
	"os"
	"os/signal"
	"time"
)

var Log *log.Logger
var Bot *bot.Bot
var Conf Config

type EchoRequestValidator struct {
	validator *validator.Validate
}

func (cv *EchoRequestValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	if err := configor.Load(&Conf, "config.yml"); err != nil {
		panic(err)
	}

	{
		logFile, err := os.OpenFile(Conf.Log.Name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
		if err != nil {
			panic(err)
		}
		Log = log.New(io.MultiWriter(logFile, os.Stderr), "[main] ", log.LstdFlags|log.Lshortfile)
	}

	dialBot()

	e := echo.New()
	e.Validator = &EchoRequestValidator{validator: validator.New()}

	e.POST("/webhook", webhookHandler)

	// Start server
	go func() {
		if err := e.Start(Conf.Server.Address); err != nil {
			e.Logger.Info("quit routine initiated...")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("cleaning up...")
	fmt.Println("releasing mirai-http session as exiting http server")
	Bot.Client.Release(Conf.MiraiHTTP.QQNumber)

	fmt.Println("shutting down http server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
