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
	GetDeviceById(ctx context.Context, id uint) (entities.Devices, error)
	UpdateDeviceStatus(ctx context.Context, device entities.Devices, deviceStatus bool) (entities.Devices, error)
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

func (d *deviceRepoService) GetDeviceById(ctx context.Context, id uint) (entities.Devices, error) {
	device, err := d.deviceRepository.FindDeviceById(ctx, id)
	return device, err
}

func (d *deviceRepoService) CreateNewDevice(ctx context.Context, device dto.RequestDeviceDto) (entities.Devices, error) {
	deviceEntity := entities.Devices{
		DeviceVendor:    device.DeviceVendor,
		DeviceName:      device.DeviceName,
		DeviceSchema:    device.DeviceSchema,
		DeviceIpAddress: device.DeviceIpAddress,
		DevicePort:      device.DevicePort,
		DeviceStatus:    false,
	}
	deviceResponse, err := d.deviceRepository.Create(ctx, deviceEntity)
	return deviceResponse, err
}

func (d *deviceRepoService) UpdateDeviceStatus(ctx context.Context, device entities.Devices, deviceStatus bool) (entities.Devices, error) {
	device.DeviceStatus = deviceStatus
	deviceResponse, err := d.deviceRepository.Update(ctx, device)
	return deviceResponse, err
}

func (d *deviceRepoService) getDevicesStrings(ctx context.Context) ([]entities.Devices, error) {
	devices, err := d.deviceRepository.FindAll(ctx)
	return devices, err
}
