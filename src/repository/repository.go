package repository

import (
	"log"
	"rsreu-class-schedule/config"
)

type ScheduleRepository interface {
	GetScheduleTypes() ([]string, error)
}

func NewRepository(cfg *config.Config) ScheduleRepository {
	switch cfg.RepositoryType {
	case config.RepositoryTypeFile:
		return newScheduleFilesRepository(cfg.SchedulesPath)
	default:
		log.Printf("Unknown repository type: %s", cfg.RepositoryType)
		return nil
	}
}
