package main

import (
	// config
	"github.com/5aradise/gather-weather/config"

	// handlers
	subscriptionHandler "github.com/5aradise/gather-weather/internal/controllers/subscription"
	weatherHandler "github.com/5aradise/gather-weather/internal/controllers/weather"
	subscriptionStorage "github.com/5aradise/gather-weather/internal/storages/subscription"

	// services
	subscriptionService "github.com/5aradise/gather-weather/internal/services/subscriber"
	validationServ "github.com/5aradise/gather-weather/internal/services/validator"
	weatherService "github.com/5aradise/gather-weather/internal/services/weather"

	// storages
	"github.com/5aradise/gather-weather/pkg/db/postgres"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var envPath = flag.String("env", "", "Path to env file")

func main() {
	flag.Parse()

	if *envPath != "" {
		err := config.Load(*envPath)
		if err != nil {
			log.Fatal("can't load env vars: ", err)
		}
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatal("can't load config: ", err)
	}

	db, err := postgres.New(postgres.Config{
		Env: cfg.Env,

		Host:     cfg.DB.Address,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Port:     cfg.DB.Port,
		Name:     cfg.DB.Name,
	})
	if err != nil {
		log.Fatal("can't init db: ", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("can't close db: ", err)
		}
	}()

	// storages
	subStor := subscriptionStorage.New(db.API())

	// services
	weatherSrv, err := weatherService.New(cfg.WeatherApiKey, sonic.Unmarshal)
	if err != nil {
		log.Fatal("can't init weather service: ", err)
	}
	validSrv := validationServ.New(weatherSrv.CheckCity)
	subSrv := subscriptionService.New(subStor, validSrv)

	// handlers
	weatherH := weatherHandler.New(weatherSrv)
	subH := subscriptionHandler.New(subSrv)

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowMethods:  []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodOptions, fiber.MethodPut, fiber.MethodDelete},
		ExposeHeaders: []string{"Link"},
	}))
	app.Use(logger.New())

	api := app.Group("/api")

	weatherH.Init(api)
	subH.Init(api)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error)
	go func() {
		serverErr <- app.Listen(net.JoinHostPort("", cfg.Server.Port))
	}()

	select {
	case s := <-interrupt:
		log.Println("signal interrupt: ", s.String())
	case err := <-serverErr:
		log.Println("server error: ", err)
	}

	err = app.Shutdown()
	if err != nil {
		log.Fatal("can't shutdown server: ", err)
	}
}
