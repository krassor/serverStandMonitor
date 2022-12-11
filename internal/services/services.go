package services

import (
	"context"

	"github.com/serverStandMonitor/internal/models/dto"
	"github.com/serverStandMonitor/internal/models/entities"
	"github.com/serverStandMonitor/internal/repositories"
)

type DevicesRepoService interface {
	GetDevices(ctx context.Context) ([]entities.Devices, error)
	CreateNewDevice(ctx context.Context, device dto.RequestDeviceDto) (entities.Devices, error)
}

type deviceRepoService struct {
	deviceRepository repositories.DevicesRepository
}

func NewdeviceRepoService(deviceRepository repositories.DevicesRepository) DevicesRepoService {
	return &deviceRepoService{
		deviceRepository: deviceRepository,
	}
}

func (d *deviceRepoService) GetDevices(ctx context.Context) ([]entities.Devices, error) {
	devices, err := d.deviceRepository.FindAll(ctx)
	return devices, err
}

func (d *deviceRepoService) CreateNewDevice(ctx context.Context, device dto.RequestDeviceDto) (entities.Devices ,error) {
	deviceEntity := entities.Devices{
		DeviceVendor:      device.DeviceVendor,
		DeviceModel:       device.DeviceModel,
		DeviceGetEndpoint: device.DeviceGetEndpoint,
		DeviceIpAddress:   device.DeviceIpAddress,
		DevicePort:        device.DevicePort,
	}
	err := d.deviceRepository.Save(ctx, deviceEntity)
	return deviceEntity, err
}
