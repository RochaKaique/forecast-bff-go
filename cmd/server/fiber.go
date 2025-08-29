package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/RochaKaique/forecastgo/internal/coordinates"
	"github.com/RochaKaique/forecastgo/internal/forecast"
	"github.com/RochaKaique/forecastgo/internal/forecast/handler"
	"github.com/RochaKaique/forecastgo/internal/forecast/service"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/spf13/viper"
)

const (
	ContextPath = "/forecast/v1"
)

type Server struct {
	app  *fiber.App
	port string
}

func (s *Server) Start() error {
	return s.app.Listen(s.port)
}

func (s *Server) Shutdown(ctx context.Context) {
	s.app.ShutdownWithContext(ctx)
}

func CreateServer(conf *viper.Viper) *Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	tr := &http.Transport{
		MaxIdleConns:          256,
		MaxIdleConnsPerHost:   64,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		// Opcional: MaxConnsPerHost: 128,
	}
	httpClient := &http.Client{
		Timeout:   10 * time.Second, // timeout total da chamada
		Transport: tr,
	}

	appName := conf.GetString("application.name")
	prometheus := fiberprometheus.New(appName)
	prometheus.RegisterAt(app, "/metrics")

	app.Use(prometheus.Middleware)
	app.Use(healthcheck.New())

	server := &Server{
		app:  app,
		port: conf.GetString("server.port"),
	}

	coordClient := coordinates.NewCooordinatesClient(httpClient, conf)
	fcClient := forecast.NewForecastClient(httpClient, conf)

	log := slog.Default()
	svc := service.New(coordClient, fcClient, log)

	api := server.app.Group(ContextPath)
	{
		api.Get("/:zipcode", handler.Register(svc))
	}

	return server
}
