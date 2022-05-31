package main

import (
	"about-go/app/default/api"
	"about-go/config"
	"about-go/pkg/controller"
	"about-go/pkg/errs"
	"about-go/pkg/logger"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

// app
func main() {
	// load env config
	cfg := config.Load()

	// set GIN_MODE
	gin.SetMode(cfg.GinMode)

	// init http server
	server := gin.New()

	// log to console&file
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(logger.Writer(logger.INFO), os.Stdout)

	// global logger
	server.Use(gin.Logger())

	// global recovery
	server.Use(gin.CustomRecoveryWithWriter(
		io.MultiWriter(logger.Writer(logger.ERROR), os.Stdout),
		errs.DefaultRecoveryHandler(),
	))

	// set trusted proxies
	if err := server.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		panic(err)
	}

	// routes group & routes
	controller.RegisterForGin(&server.RouterGroup, &api.UserController{})

	// listen port
	server.Run(":80")
}
