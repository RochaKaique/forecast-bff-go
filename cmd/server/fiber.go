package server

import (
	"context"
	"fmt"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/spf13/viper"
)

type Server struct {
	app  *fiber.App
	port string
}

func (s Server) Start() error {
	return s.app.Listen(s.port)
}

func (s Server) Shutdown(ctx context.Context) {
	s.app.ShutdownWithContext(ctx)
}

func CreateServer(conf *viper.Viper) Server {
	app := fiber.New()

	appName := conf.GetString("application.name")
	prometheus := fiberprometheus.New(appName)
	prometheus.RegisterAt(app, "/metrics")

	app.Use(prometheus.Middleware)
	app.Use(healthcheck.New())

	server := Server{
		app: app,
		port: conf.GetString("server.port"),
	}

	fmt.Println(appName)
	return server
}
