package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/serverStandMonitor/internal/models/dto"
	"github.com/serverStandMonitor/internal/models/entities"
	"github.com/serverStandMonitor/internal/services"
	"github.com/serverStandMonitor/pkg/utils"
	"net/http"
	"net/url"
)

type DeviceHandlers interface {
	CreateDevice(w http.ResponseWriter, r *http.Request)
}

type deviceHandler struct {
	deviceService services.DevicesRepoService
}

func NewDeviceHandler(deviceService services.DevicesRepoService) DeviceHandlers {
	return &deviceHandler{
		deviceService: deviceService,
	}
}

func (d *deviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	deviceDto := dto.RequestDeviceDto{}
	err := json.NewDecoder(r.Body).Decode(&deviceDto)
	if err != nil {
		log.Warn().Msgf("Error decode json: %s", err)
		utils.Err(w, http.StatusInternalServerError, err)
		return
	}
	deviceUrl := fmt.Sprintf("%s://%s:%s", deviceDto.DeviceSchema, deviceDto.DeviceIpAddress, deviceDto.DevicePort)
	_, err = url.Parse(deviceUrl)
	if err != nil {
		log.Warn().Msgf("Error parse URL: %s", err)
		utils.Err(w, http.StatusInternalServerError, err)
		return
	}
	deviceEntity := entities.Devices{}
	deviceEntity, err = d.deviceService.CreateNewDevice(r.Context(), deviceDto)

	if err != nil {
		log.Error().Msgf("Error creating device: %s", err)
		utils.Err(w, http.StatusInternalServerError, err)
		return
	}

	responseDeviceParams := dto.ResponseDeviceParams{}
	responseDeviceParams.DeviceId = deviceEntity.ID
	responseDevice := utils.Message(true, responseDeviceParams)
	utils.Json(w, http.StatusOK, responseDevice)
}
