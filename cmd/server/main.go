package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"wallet-service/internal/app"
	"wallet-service/internal/handler/api_1_wallet_get"
	"wallet-service/internal/handler/api_1_wallet_post"
	"wallet-service/internal/service/wallet_service"
	wallet_storage "wallet-service/internal/storage/postgres/wallet"
)

func main() {
	config := app.MustSetUpConfig()
	app.MakeMigrations(true, config)

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbClient := app.MustSetUpDb("postgres", config)

	e := echo.New()

	walletStorage := wallet_storage.New(dbClient, log)

	walletService := wallet_service.New(log, walletStorage)

	apiWalletPost := api_1_wallet_post.New(log, walletService)
	apiWalletGet := api_1_wallet_get.New(log, walletService)

	e.POST("/api/v1/wallet", apiWalletPost.Handle)
	e.GET("/api/v1/wallets/:WALLET_UUID", apiWalletGet.Handle)

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", config.Server.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", "error", err)
		}
	}()

	<-stop
	log.Info("received shutdown signal, starting shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error("failed to gracefully shut down server", "error", err)
	}

	log.Info("server gracefully stopped")

}
