package dto

type RequestDeviceDto struct {
	DeviceVendor      string `json:"deviceVendor"`
	DeviceModel       string `json:"deviceModel"`
	DeviceGetEndpoint string `json:"deviceGetEndpoint"`
	DeviceIpAddress   string `json:"deviceIpAddress"`
	DevicePort        string `json:"devicePort"`
}
