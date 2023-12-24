package repository

import (
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"rsreu-class-schedule/config"
	"strings"
)

type ScheduleFilesRepository struct {
	schedulePath string
}

func newScheduleFilesRepository(schedulePath string) *ScheduleFilesRepository {
	return &ScheduleFilesRepository{schedulePath: schedulePath}
}

func (r *ScheduleFilesRepository) GetScheduleTypes() ([]string, error) {
	scheduleTypes := make([]string, 0)

	err := filepath.WalkDir(r.schedulePath, func(path string, d fs.DirEntry, err error) error {
		name := d.Name()

		if name == config.SchedulesPath || !d.IsDir() {
			return nil
		}

		for _, firstLevelDir := range scheduleTypes {
			if strings.Contains(path, firstLevelDir) {
				return nil
			}
		}

		scheduleTypes = append(scheduleTypes, name)

		return nil
	})
	if err != nil {
		log.Printf("failed get schedule types form %s: %v", r.schedulePath, err)
		return nil, err
	}

	log.Println("Types of schedules are obtained: ", scheduleTypes)

	return scheduleTypes, nil
}

func (r *ScheduleFilesRepository) GetScheduleFile(faculty string, studyType StudyType) (string, error) {
	pathToFile := path.Join(r.schedulePath, faculty, studyType.String())

	var name string

	err := filepath.Walk(pathToFile, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		name = info.Name()

		return nil
	})
	if err != nil {
		log.Printf("failed get schedule types form %s: %v", r.schedulePath, err)
		return "", err
	}

	return path.Join(pathToFile, name), nil
}
