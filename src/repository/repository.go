package repository

import (
	"log"
	"rsreu-class-schedule/config"

	"github.com/xuri/excelize/v2"
)

type ScheduleRepository interface {
	GetScheduleTypes() ([]string, error)
	GetScheduleFile(faculty string, studyType StudyType) (*excelize.File, error)
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
