package dto

type ResponseDeviceDto struct {
	DeviceVendor      string `json:"deviceVendor"`
	DeviceModel       string `json:"deviceModel"`
	DeviceGetEndpoint string `json:"deviceGetEndpoint"`
	DeviceIpAddress   string `json:"deviceIpAddress"`
	DevicePort        string `json:"devicePort"`
}

type ResponseDeviceParams struct {
	DeviceId uint `json:"deviceId"`
}
