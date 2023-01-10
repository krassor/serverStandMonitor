package entities

import (
	"gorm.io/gorm"
)

type Devices struct {
	gorm.Model
	DeviceVendor    string `gorm:"column:deviceVendor"`
	DeviceName      string `gorm:"column:deviceName"`
	DeviceSchema    string `gorm:"column:deviceSchema"`
	DeviceIpAddress string `gorm:"column:deviceIpAddress"`
	DevicePort      string `gorm:"column:devicePort"`
	DeviceStatus    bool   `gorm:"column:deviceStatus;default:false"`
}
