package repository

//go:generate go-enum -f=$GOFILE
/*
ENUM(
	Бакалавриат
	Магистратура
	Аспирантура
	Специалитет
)
*/
type StudyType string
