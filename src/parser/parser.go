package parser

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"rsreu-class-schedule/entities"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

var (
	startTable = "дни"

	regexLesson = regexp.MustCompile(`(?P<name>.*\n)(?P<teacher>.*?\s)(?P<classroom>\d+.*)`)

	nilLessonError = errors.New("lesson is empty")

	regexGroupName      = "name"
	regexGroupTeacher   = "teacher"
	regexGroupClassroom = "classroom"
	regexGroupDates     = "dates"

	countParsedElements          = 4
	countParsedElementsWithDates = 5
)

type axisCounter struct {
	coll    []int
	row     int
	minColl []int
}

func newAxisCounter() *axisCounter {
	return &axisCounter{
		coll:    []int{65},
		row:     1,
		minColl: []int{65},
	}
}

func (a *axisCounter) nextColl() {
	for i := range a.coll {
		if a.coll[i] == 90 {
			if i == len(a.coll)-1 {
				a.coll = append(a.coll, 65)
				break
			}

			continue
		}

		a.coll[i]++
	}
}

func (a *axisCounter) setMinColl() {
	copy(a.minColl, a.coll)
}

func (a *axisCounter) nextRow() {
	a.row++
	copy(a.coll, a.minColl)
}

func (a *axisCounter) String() string {
	colls := make([]string, len(a.coll))

	for i := range a.coll {
		colls[i] = string(a.coll[i])
	}

	return fmt.Sprintf("%s%d", strings.Join(colls, ""), a.row)
}

func ParseSchedule(file *excelize.File) ([]*entities.Schedule, error) {
	var err error

	axis := newAxisCounter()

	if len(file.GetSheetList()) == 0 {
		err := fmt.Errorf("not found sheets")
		log.Printf("failed open file %s: %v", file.Path, err)
		return nil, err
	}

	sheet := file.GetSheetList()[0]

	rewindAxisToStart(file, sheet, axis)
	schedules := createSchedules(file, sheet, axis)

	for {
		weekday, err := getWeekday(file, sheet, axis)
		if err != nil {
			break
		}
		lessonTime, err := getLessonTime(file, sheet, axis)
		if err != nil {
			log.Printf("failed get lesson time %s: %v", file.Path, err)
			break
		}
		weekType, err := getWeekType(file, sheet, axis)
		if err != nil {
			log.Printf("failed get week type %s: %v", file.Path, err)
			break
		}

		for _, schedule := range schedules {
			name, teacher, classroom, dates, err := getLessonInfo(file, sheet, axis)
			if err != nil {
				if err == nilLessonError {
					continue
				}
				log.Printf("failed get lesson info %s: %v", file.Path, err)
				break
			}

			schedule.SetLesson(weekday, weekType, entities.NewLesson(name, teacher, classroom, lessonTime, dates))
		}

		axis.nextRow()
	}

	return schedules, err
}

func getLessonInfo(
	file *excelize.File,
	sheet string,
	axis *axisCounter,
) (
	name, teacher, classroom, dates string,
	err error,
) {
	defer axis.nextColl()

	value, err := file.GetCellValue(sheet, axis.String())
	if err != nil {
		log.Printf("failed get value from %s: %v", axis.String(), err)
		return "", "", "", "", err
	}

	matches := regexLesson.FindStringSubmatch(value)

	if len(matches) < countParsedElements {
		if value == "" {
			return "", "", "", "", nilLessonError
		}
		return strings.TrimSpace(value), "", "", "", nil
	}

	name = strings.TrimSpace(matches[regexLesson.SubexpIndex(regexGroupName)])
	teacher = strings.TrimSpace(matches[regexLesson.SubexpIndex(regexGroupTeacher)])
	classroom = strings.TrimSpace(matches[regexLesson.SubexpIndex(regexGroupClassroom)])

	if len(matches) == countParsedElementsWithDates {
		dates = strings.TrimSpace(matches[regexLesson.SubexpIndex(regexGroupDates)])
	}

	return
}

func getWeekType(file *excelize.File, sheet string, axis *axisCounter) (entities.WeekType, error) {
	defer axis.nextColl()

	value, err := file.GetCellValue(sheet, axis.String())
	if err != nil {
		log.Printf("failed get value from %s: %v", axis.String(), err)
		return 0, err
	}

	weekType := strings.ToLower(strings.TrimSpace(value))

	switch weekType {
	case "числ.":
		return entities.WeekTypeNumerator, nil
	case "знам.":
		return entities.WeekTypeDenominator, nil
	default:
		return 0, fmt.Errorf("unexpected weektype: %s", weekType)
	}
}

func getLessonTime(file *excelize.File, sheet string, axis *axisCounter) (string, error) {
	defer axis.nextColl()

	value, err := file.GetCellValue(sheet, axis.String())
	if err != nil {
		log.Printf("failed get value from %s: %v", axis.String(), err)
		return "", err
	}

	return strings.ToLower(strings.TrimSpace(value)), nil
}

func getWeekday(file *excelize.File, sheet string, axis *axisCounter) (time.Weekday, error) {
	defer axis.nextColl()

	value, err := file.GetCellValue(sheet, axis.String())
	if err != nil {
		log.Printf("failed get value from %s: %v", axis.String(), err)
		return 0, err
	}

	weekday := strings.ToLower(strings.TrimSpace(value))

	switch weekday {
	case "понедельник":
		return time.Monday, nil
	case "вторник":
		return time.Tuesday, nil
	case "среда":
		return time.Wednesday, nil
	case "четверг":
		return time.Thursday, nil
	case "пятница":
		return time.Friday, nil
	case "суббота":
		return time.Saturday, nil
	case "воскресенье":
		return time.Sunday, nil
	default:
		return 0, fmt.Errorf("unexpected weekday: %s", weekday)
	}
}

func createSchedules(file *excelize.File, sheet string, axis *axisCounter) (schedules []*entities.Schedule) {
	for {
		value, err := file.GetCellValue(sheet, axis.String())
		if err != nil {
			log.Printf("failed get value from %s: %v", axis.String(), err)
			break
		}

		if value == "" {
			break
		}

		schedules = append(schedules, entities.NewSchedule(strings.TrimSpace(value)))

		axis.nextColl()
	}

	axis.nextRow()

	return schedules
}

func rewindAxisToStart(file *excelize.File, sheet string, axis *axisCounter) {
	for i := 0; i < 100; i++ {
		for j := 0; j < 2; j++ {
			value, _ := file.GetCellValue(sheet, axis.String())

			if strings.ToLower(strings.TrimSpace(value)) == startTable {
				axis.setMinColl()

				// скипаем бесполезные ячейки
				axis.nextColl()
				axis.nextColl()
				axis.nextColl()
				return
			}

			axis.nextColl()
		}

		axis.nextRow()
	}
}
