package schedule

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rsreu-class-schedule/server/errors"
)

func (c *Controller) GetScheduleTypes(context *gin.Context) {
	scheduleTypes, err := c.repository.GetScheduleTypes()
	if err != nil {
		context.IndentedJSON(errors.NewError500(err, "failed get schedule types"))
		return
	}

	context.IndentedJSON(http.StatusOK, scheduleTypes)
}