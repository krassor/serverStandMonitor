package repositories

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/serverStandMonitor/internal/models/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DevicesRepository interface {
	FindAll(ctx context.Context) ([]entities.Devices, error)
	Save(ctx context.Context, device entities.Devices) error
}

type deviceRepository struct {
	DB *gorm.DB
}

func NewDevicesRepository() DevicesRepository {
	username := os.Getenv("DEVICES_DB_USER")
	password := os.Getenv("DEVICES_DB_PASSWORD")
	dbName := os.Getenv("DEVICES_DB_NAME")
	dbHost := os.Getenv("DEVICES_DB_HOST")
	dbPort := os.Getenv("DEVICES_DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, username, dbName, password)
	fmt.Println(dsn)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error().Msgf("Error gorm.Open(): %s", err)
	}
	log.Info().Msg("gorm have connected to database")

	err = conn.Debug().AutoMigrate(&entities.Devices{}) //Миграция базы данных
	if err != nil {
		log.Error().Msgf("Error gorm.AutoMigrate(): %s", err)
	}
	log.Info().Msg("gorm have connected to database")

	return &deviceRepository{
		DB: conn,
	}
}

func (d *deviceRepository) FindAll(ctx context.Context) ([]entities.Devices, error) {
	var devices []entities.Devices
	tx := d.DB.WithContext(ctx).Find(&devices)
	if tx.Error != nil {
		//log.Warn().Msgf("Error finding devices in repo: %s", tx.Error)
		return []entities.Devices{}, tx.Error
	}

	return devices, nil
}

func (d *deviceRepository) Save(ctx context.Context, device entities.Devices) error {
	tx := d.DB.WithContext(ctx).Save(&device)
	if tx.Error != nil {
		//log.Warn().Msgf("Error saving device in repo: %s", tx.Error)
		return tx.Error
	}
	return nil
}
