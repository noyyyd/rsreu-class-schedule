package entities

import "time"

//go:generate go-enum -f=$GOFILE
/*
ENUM(
	numerator,
	denominator
)
*/
type WeekType int

type Schedule struct {
	Group     string                 `json:"group"`
	Monday    map[WeekType][]*Lesson `json:"monday,omitempty"`
	Tuesday   map[WeekType][]*Lesson `json:"tuesday,omitempty"`
	Wednesday map[WeekType][]*Lesson `json:"wednesday,omitempty"`
	Thursday  map[WeekType][]*Lesson `json:"thursday,omitempty"`
	Friday    map[WeekType][]*Lesson `json:"friday,omitempty"`
	Saturday  map[WeekType][]*Lesson `json:"saturday,omitempty"`
	Sunday    map[WeekType][]*Lesson `json:"sunday,omitempty"`
}

func NewSchedule(group string) *Schedule {
	return &Schedule{
		Group:     group,
		Monday:    make(map[WeekType][]*Lesson),
		Tuesday:   make(map[WeekType][]*Lesson),
		Wednesday: make(map[WeekType][]*Lesson),
		Thursday:  make(map[WeekType][]*Lesson),
		Friday:    make(map[WeekType][]*Lesson),
		Saturday:  make(map[WeekType][]*Lesson),
		Sunday:    make(map[WeekType][]*Lesson),
	}
}

func (s *Schedule) SetLesson(weekday time.Weekday, weekType WeekType, lesson *Lesson) {
	switch weekday {
	case time.Monday:
		s.Monday[weekType] = append(s.Monday[weekType], lesson)
	case time.Tuesday:
		s.Tuesday[weekType] = append(s.Tuesday[weekType], lesson)
	case time.Wednesday:
		s.Wednesday[weekType] = append(s.Wednesday[weekType], lesson)
	case time.Thursday:
		s.Thursday[weekType] = append(s.Thursday[weekType], lesson)
	case time.Friday:
		s.Friday[weekType] = append(s.Friday[weekType], lesson)
	case time.Saturday:
		s.Saturday[weekType] = append(s.Saturday[weekType], lesson)
	case time.Sunday:
		s.Sunday[weekType] = append(s.Sunday[weekType], lesson)
	}
}

type Lesson struct {
	Name      string `json:"name"`
	Teacher   string `json:"teacher"`
	Classroom string `json:"classroom"`
	Time      string `json:"time"`
	Dates     string `json:"dates,omitempty"`
}

func NewLesson(name string, teacher string, classroom string, time string, dates string) *Lesson {
	return &Lesson{Name: name, Teacher: teacher, Classroom: classroom, Time: time, Dates: dates}
}
