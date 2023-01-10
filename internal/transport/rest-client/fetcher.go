package rest_client

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/serverStandMonitor/internal/models/entities"
	"github.com/serverStandMonitor/internal/services"
	"github.com/serverStandMonitor/internal/transport/rest-client/client"
)

const (
	fetcherDuration      time.Duration = 5
	timeoutGetDeviceList time.Duration = 3
)

type DeviceFetcher struct {
	client  client.DeviceStatusClient
	service services.DevicesRepoService
	exit    chan struct{}
}

func NewDeviceFetcher(service services.DevicesRepoService) *DeviceFetcher {
	return &DeviceFetcher{client: client.NewDefaultDevice(&http.Client{}), service: service, exit: make(chan struct{})}
}

func (f *DeviceFetcher) Start(ctx context.Context) {

	var wg sync.WaitGroup
	for {
		select {
		case <-f.exit:
			log.Info().Msgf("Exit device fetcher")
			return
		default:
		}
		entityList := f.getDeviceList(ctx, timeoutGetDeviceList)
		log.Info().Msgf("select default")

		wg.Add(len(entityList))
		log.Info().Msgf("workgroup add %d elements", len(entityList))

		for _, e := range entityList {
			go func(entity entities.Devices) {
				defer wg.Done()
				log.Info().Msgf("Enter anonymous go routine")

				deviceUrl := fmt.Sprintf("%s://%s:%s", entity.DeviceSchema, entity.DeviceIpAddress, entity.DevicePort)
				log.Info().Msgf("device URL: %s", deviceUrl)

				status, err := f.client.GetStatus(ctx, deviceUrl)
				if err != nil {
					log.Error().Msgf("Error get device %s %s %s status in fetcher.Start(): %s", entity.ID, entity.DeviceVendor, entity.DeviceName, err)
				}
				log.Info().Msgf("Device status: %t", status)

				_, err = f.service.UpdateDeviceStatus(ctx, entity, status)
				if err != nil {
					log.Error().Msgf("Error update device %s %s %s status in fetcher.Start(): %s", entity.ID, entity.DeviceVendor, entity.DeviceName, err)
				}
				log.Info().Msgf("End of anonymous go routine")
			}(e)
		}
		wg.Wait()
		time.Sleep(fetcherDuration * time.Second)

	}
}

func (f *DeviceFetcher) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("error shutdown device fetcher: %s", ctx.Err())
		default:
			f.exit <- struct{}{}
			return nil
		}
	}
}

func (f *DeviceFetcher) getDeviceList(ctx context.Context, timeoutDuration time.Duration) []entities.Devices {
	var entityList []entities.Devices
	var err error
	for {
		entityList, err = f.service.GetDevices(ctx)
		if err != nil {
			log.Error().Msgf("Fetcher: cannot get data from repo: %s", err)
			time.Sleep(timeoutDuration)
			continue
		}

		break
	}
	return entityList
}
