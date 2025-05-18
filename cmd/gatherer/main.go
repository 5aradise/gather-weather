package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/5aradise/gather-weather/config"
	"github.com/5aradise/gather-weather/pkg/db/postgres"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
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

	app.Get("", func(c fiber.Ctx) error {
		var tables []string
		err = db.API().WithContext(c.Context()).Raw(`
		SELECT tablename
		FROM pg_catalog.pg_tables
		WHERE schemaname NOT IN ('pg_catalog', 'information_schema');
		`).Scan(&tables).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(tables)
	})

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
