package rest_client

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/serverStandMonitor/internal/models/entities"
	"github.com/serverStandMonitor/internal/services"
	"github.com/serverStandMonitor/internal/transport/rest-client/client"
	"net/http"
	"sync"
	"time"
)

const (
	fetcherDuration time.Duration = 5
)

type DeviceFetcher struct {
	client  client.DeviceStatusClient
	service services.DevicesRepoService
	exitCh  chan bool
}

func NewDeviceFetcher(service services.DevicesRepoService) *DeviceFetcher {
	return &DeviceFetcher{client: client.NewDefaultDevice(&http.Client{}), service: service, exitCh: make(chan bool)}
}

func (f *DeviceFetcher) Start(ctx context.Context) {
	var entitySlice []entities.Devices
	for {
		ent, err := f.service.GetDevices(ctx)
		if err != nil {
			log.Error().Msgf("Fetcher: cannot get data from repo: %s", err)
			time.Sleep(3)
			continue
		}
		entitySlice = ent
		break
	}

	log.Info().Msgf("entitySlice: %v", entitySlice)

	var wg sync.WaitGroup
	for {
		select {
		case <-f.exitCh:
			log.Info().Msgf("<-f.exitCh")
			return
		default:
		}
		log.Info().Msgf("select default")
		wg.Add(len(entitySlice))
		log.Info().Msgf("workgroup add %d elements", len(entitySlice))
		for i, e := range entitySlice {
			go func(entity entities.Devices) {
				defer wg.Done()
				log.Info().Msgf("%d Enter anonymous go routine", i)
				deviceUrl := fmt.Sprintf("%s://%s:%s", entity.DeviceSchema, entity.DeviceIpAddress, entity.DevicePort)
				log.Info().Msgf("%d device URL: %s", i, deviceUrl)
				status, err := f.client.GetStatus(ctx, deviceUrl)
				if err != nil {
					log.Error().Msgf("%d Error get device %s %s %s status in fetcher.Start(): %s", i, entity.ID, entity.DeviceVendor, entity.DeviceName, err)
				}
				log.Info().Msgf("%d Device status: %b", i, status)
				_, err = f.service.UpdateDeviceStatus(ctx, entity, status)
				if err != nil {
					log.Error().Msgf("%d Error update device %s %s %s status in fetcher.Start(): %s", i, entity.ID, entity.DeviceVendor, entity.DeviceName, err)
				}
				log.Info().Msgf("%d End of anonymous go routine", i)
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
			return fmt.Errorf("Error shutdown device fetcher: %s", ctx.Err())
		default:
			f.exitCh <- true
		}
	}
}
