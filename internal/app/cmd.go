package app

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `yml:"host"`
		Port string `yml:"port"`
	}
	DB struct {
		Host string `yml:"host"`
		Port string `yml:"port"`
		User string `yml:"user"`
		Pass string `yml:"pass"`
		Name string `yml:"name"`
	}
}

func MustSetUpConfig() *Config {
	var appConfig Config
	viper.SetConfigName("app")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&appConfig); err != nil {
		panic(err)
	}

	return &appConfig
}

func MustSetUpDb(driverName string, config *Config) *sqlx.DB {
	dbInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Pass,
		config.DB.Name,
	)
	db, err := sqlx.Connect(driverName, dbInfo)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(70)
	db.SetMaxIdleConns(35)

	return db
}

func MakeMigrations(up bool, config *Config) {
	dbLine := fmt.Sprintf("postgres://%s:%s@db:%s/%s?sslmode=disable",
		config.DB.User,
		config.DB.Pass,
		config.DB.Port,
		config.DB.Name,
	)
	m, err := migrate.New("file://migrations", dbLine)
	if err != nil {
		panic(err)
	}

	if up {
		err = m.Up()
		if err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				panic(err)
			}
		}
	} else {
		err = m.Down()
		if !errors.Is(err, migrate.ErrNoChange) {
			panic(err)
		}
	}
}
