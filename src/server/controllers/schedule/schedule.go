package schedule

import (
	"fmt"
	"log"
	"net/http"
	"rsreu-class-schedule/parser"
	"rsreu-class-schedule/repository"
	"rsreu-class-schedule/server/errors"

	"github.com/gin-gonic/gin"
)

type GetRequest struct {
	Faculty string               `json:"faculty"`
	Type    repository.StudyType `json:"type"`
}

func (r *GetRequest) ValidateRequest() error {
	if r.Faculty == "" {
		return fmt.Errorf("empty faculty")
	}
	if r.Type == "" {
		return fmt.Errorf("empty schedule type")
	}
	return nil
}

func (c *Controller) GetSchedule(context *gin.Context) {
	req := new(GetRequest)

	if err := c.unmarshalRequest(context, req); err != nil {
		context.IndentedJSON(errors.NewError400(err, "failed get request"))
		return
	}

	file, err := c.repository.GetScheduleFile(req.Faculty, req.Type)
	if err != nil {
		context.IndentedJSON(errors.NewError500(err, "failed get file"))
		return
	}

	schedules, err := parser.ParseSchedule(file)
	if err != nil {
		log.Printf("failed parse schedule %s: %v", file.Path, err)
		context.IndentedJSON(errors.NewError500(err, ""))
		return
	}

	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Headers", "Content-Type")
	context.IndentedJSON(http.StatusOK, schedules)
}
