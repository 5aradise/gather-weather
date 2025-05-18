package postgres

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	Database struct {
		conn *gorm.DB
	}

	Config struct {
		Env string

		Host, User, Password, Port, Name string
	}
)

func New(cfg Config) (*Database, error) {
	logLvl := logger.Silent
	switch cfg.Env {
	case "dev":
		logLvl = logger.Info
	case "prod":
		logLvl = logger.Error
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLvl),
	})
	if err != nil {
		return nil, err
	}

	sql, err := conn.DB()
	if err != nil {
		return nil, err
	}
	if err = sql.Ping(); err != nil {
		return nil, err
	}

	return &Database{conn}, nil
}

func (db *Database) API() *gorm.DB {
	return db.conn
}

func (db *Database) Close() error {
	if db.conn == nil {
		return errors.New("db connection is already closed")
	}
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
