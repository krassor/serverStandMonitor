package main

import (
	"context"
	"time"

	"github.com/serverStandMonitor/internal/graceful"
	"github.com/serverStandMonitor/internal/logger"
	"github.com/serverStandMonitor/internal/repositories"
	"github.com/serverStandMonitor/internal/services"
	httpServer "github.com/serverStandMonitor/internal/transport/rest"
	"github.com/serverStandMonitor/internal/transport/rest/handlers"
	"github.com/serverStandMonitor/internal/transport/rest/routers"
)

func main() {
	logger.InitLogger()

	deviceRepository := repositories.NewDevicesRepository()
	deviceRepoService := services.NewdeviceRepoService(deviceRepository)
	deviceHandler := handlers.NewDeviceHandler(deviceRepoService)
	deviceRouter := routers.NewDeviceRouter(deviceHandler)
	deviceHttpServer := httpServer.NewHttpServer(deviceRouter)

	maxSecond := 10 * time.Second
	graceful.GracefulShutdown(
		context.Background(),
		maxSecond,
		map[string]graceful.Operation{
			"http": func(ctx context.Context) error {
				return deviceHttpServer.Shutdown(ctx)
			},
		},
	)

	deviceHttpServer.Listen()
}
