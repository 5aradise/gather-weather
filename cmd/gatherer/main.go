package main

import (
	// config
	"net/http"

	root "github.com/5aradise/gather-weather"
	"github.com/5aradise/gather-weather/config"

	// handlers
	subscriptionHandler "github.com/5aradise/gather-weather/internal/controllers/subscription"
	weatherHandler "github.com/5aradise/gather-weather/internal/controllers/weather"
	subscriptionStorage "github.com/5aradise/gather-weather/internal/storages/subscription"

	// services
	mailService "github.com/5aradise/gather-weather/internal/services/mailer"
	subscriptionService "github.com/5aradise/gather-weather/internal/services/subscriber"
	validationServ "github.com/5aradise/gather-weather/internal/services/validator"
	weatherService "github.com/5aradise/gather-weather/internal/services/weather"

	// storages
	"github.com/5aradise/gather-weather/pkg/db/postgres"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
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

	if err := postgres.Migrate(db.API()); err != nil {
		log.Fatal("can't migrate db: ", err)
	}

	// storages
	subStor := subscriptionStorage.New(db.API())

	// services
	mailSrv := mailService.New(cfg.Mail.Host, cfg.Mail.Port, cfg.Mail.Sender, cfg.Mail.Password)
	weatherSrv, err := weatherService.New(cfg.WeatherApiKey, sonic.Unmarshal)
	if err != nil {
		log.Fatal("can't init weather service: ", err)
	}
	validSrv := validationServ.New(weatherSrv.CheckCity)
	subSrv, err := subscriptionService.New(subStor, validSrv)
	if err != nil {
		log.Fatal("can't init subscription service: ", err)
	}

	// handlers
	weatherH := weatherHandler.New(weatherSrv)
	subH := subscriptionHandler.New(subSrv, mailSrv, weatherSrv)

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	prometheus := fiberprometheus.New("my-service-name")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowMethods:  strings.Join([]string{fiber.MethodGet, fiber.MethodPost, fiber.MethodOptions, fiber.MethodPut, fiber.MethodDelete}, ","),
		ExposeHeaders: "Link",
	}))
	app.Use(logger.New())

	app.Use(filesystem.New(filesystem.Config{
		Root: http.FS(root.Public),

		Index: "public/index.html",

		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api") || strings.HasPrefix(c.Path(), "/metrics")
		},
	}))

	api := app.Group("/api")

	weatherH.Init(api)
	subH.Init(api)

	// subH.RunMailing()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error)
	go func() {
		log.Printf("server is running on port %s\n", cfg.Server.Port)
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
