package main

import (
	"context"
	"time"

	"github.com/serverStandMonitor/internal/graceful"
	"github.com/serverStandMonitor/internal/logger"
	"github.com/serverStandMonitor/internal/repositories"
	"github.com/serverStandMonitor/internal/services"
	telegramBot "github.com/serverStandMonitor/internal/telegram"
	httpServer "github.com/serverStandMonitor/internal/transport/rest-sever"
	"github.com/serverStandMonitor/internal/transport/rest-sever/handlers"
	"github.com/serverStandMonitor/internal/transport/rest-sever/routers"

	fetcher "github.com/serverStandMonitor/internal/transport/rest-client"
)

func main() {
	logger.InitLogger()

	deviceRepository := repositories.NewDevicesRepository()
	deviceRepoService := services.NewdeviceRepoService(deviceRepository)
	deviceHandler := handlers.NewDeviceHandler(deviceRepoService)
	deviceRouter := routers.NewDeviceRouter(deviceHandler)
	deviceHttpServer := httpServer.NewHttpServer(deviceRouter)

	deviceTgBot := telegramBot.NewBot(deviceRepoService)

	fetcher := fetcher.NewDeviceFetcher(deviceRepoService)

	maxSecond := 15 * time.Second
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	graceful.GracefulShutdown(
		ctx,
		maxSecond,
		map[string]graceful.Operation{
			"http": func(ctx context.Context) error {
				return deviceHttpServer.Shutdown(ctx)
			},
			"tgBot": func(ctx context.Context) error {
				return deviceTgBot.Shutdown(ctx)
			},
		},
	)

	go deviceHttpServer.Listen()
	go deviceTgBot.Update(ctx, 60)
	fetcher.Start(ctx)

}
