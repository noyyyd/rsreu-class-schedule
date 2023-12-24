package parser

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"rsreu-class-schedule/entities"
	"strings"
	"time"
)

const (
	startTable        = "дни"
	skipCellsForGroup = 3

	regexGroupName      = "name"
	regexGroupTeacher   = "teacher"
	regexGroupClassroom = "classroom"
	regexGroupDates     = "dates"

	countParsedElements          = 4
	countParsedElementsWithDates = 5
)

func ParseSchedule(filePath string) ([]*entities.Schedule, error) {
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		log.Printf("failed open file %s: %v", filePath, err)
		return nil, err
	}

	if len(file.Sheets) == 0 {
		err = fmt.Errorf("xlsx file %s doesn't have sheets: %v", filePath, err)
		log.Println(err)
		return nil, err
	}

	rows := file.Sheets[0].Rows

	rowStartID, cellStartID := getStartRowID(file)
	schedules := createSchedulesForGroups(rows, rowStartID, cellStartID)

	countGroup := len(schedules)

	// так как уже прочитали данные групп
	rowStartID += 1

	var oldWeekday time.Weekday

	cellStartIDIndex := cellStartID

	// будем парсить до момента пока не словим конец, то есть когда день недели и время пары не будет пустое
	for rows[rowStartID].Cells[cellStartIDIndex].Value != rows[rowStartID].Cells[cellStartIDIndex+1].Value {
		weekday, err := getWeekday1(rows[rowStartID].Cells[cellStartIDIndex].Value)
		if err != nil {
			weekday = oldWeekday
		}

		if oldWeekday != weekday {
			oldWeekday = weekday
		}

		// так как уже прочитали день недели
		cellStartIDIndex += 1

		var lessonsTime string
		var cells []*xlsx.Cell

		// считываем данные только по одной паре, поэтому добавляем двойку чтобы не читать лишнего
		for i, row := range rows[rowStartID : rowStartID+2] {
			cells = row.Cells

			if i == 0 {
				lessonsTime = getLessonsTime(cells[cellStartIDIndex].Value)
			}

			// делаем +1 так как уже считали время
			cells = cells[cellStartIDIndex+1:]

			var weekType entities.WeekType

			// за одну итерацию берём значения для всех групп для одного типа недели
			// +1 делаем чтобы захватить тип недели(числитель или знаменатель)
			for i, cell := range cells[:countGroup+1] {
				if i == 0 {
					weekType, err = getWeekType1(cell.Value)
					if err != nil {
						log.Printf("failed parse lessons: %v", err)
						break
					}
					continue
				}

				if cell.Value == "" {
					continue
				}

				lesson, err := parseLesson(cell.Value, lessonsTime)
				if err != nil {
					log.Println(err)
					break
				}

				schedules[i-1].SetLesson(weekday, weekType, lesson)
			}
		}

		rowStartID += 2
		cellStartIDIndex = cellStartID
	}

	return schedules, nil
}

func getStartRowID(file *xlsx.File) (int, int) {
	for iRow, row := range file.Sheets[0].Rows {
		for iCell, cell := range row.Cells {
			if strings.ToLower(cell.Value) == startTable {
				return iRow, iCell
			}
		}
	}
	return 0, 0
}

func createSchedulesForGroups(rows []*xlsx.Row, rowStartID, cellStartID int) (schedules []*entities.Schedule) {
	cells := rows[rowStartID].Cells

	// добавляем skipCellsForGroup чтобы пропустить ячейки "дни", "часы" и
	// одну пустую для указания числителя и знаменателя
	for i := cellStartID + skipCellsForGroup; i < len(cells); i++ {
		if cells[i].Value != "" {
			schedules = append(schedules, entities.NewSchedule(strings.TrimSpace(cells[i].Value)))
		}
	}

	return schedules
}

func getWeekday1(value string) (time.Weekday, error) {
	if value == "" {
		return 0, fmt.Errorf("empty weekday")
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

func getLessonsTime(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func getWeekType1(value string) (entities.WeekType, error) {
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

func parseLesson(value string, lessonsTime string) (*entities.Lesson, error) {
	matches := regexLesson.FindStringSubmatch(value)

	if len(matches) < countParsedElements {
		return entities.NewLesson(strings.TrimSpace(value), "", "", lessonsTime, ""), nil
	}

	name := strings.TrimSpace(matches[regexLesson.SubexpIndex(regexGroupName)])
	teacher := strings.TrimSpace(matches[regexLesson.SubexpIndex(regexGroupTeacher)])
	classroom := strings.TrimSpace(matches[regexLesson.SubexpIndex(regexGroupClassroom)])

	var dates string

	if len(matches) == countParsedElementsWithDates {
		dates = strings.TrimSpace(matches[regexLesson.SubexpIndex(regexGroupDates)])
	}

	return entities.NewLesson(name, teacher, classroom, lessonsTime, dates), nil
}
