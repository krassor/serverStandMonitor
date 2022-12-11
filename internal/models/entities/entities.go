package entities

import (
	"gorm.io/gorm"
)

type Devices struct {
	gorm.Model
	DeviceVendor      string `gorm:"column:deviceVendor"`
	DeviceModel       string `gorm:"column:deviceModel"`
	DeviceGetEndpoint string `gorm:"column:deviceGetEndpoint"`
	DeviceIpAddress   string `gorm:"column:deviceIpAddress"`
	DevicePort        string `gorm:"column:devicePort"`
}
