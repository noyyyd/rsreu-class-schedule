package schedule

import (
	"net/http"
	"rsreu-class-schedule/server/errors"

	"github.com/gin-gonic/gin"
)

func (c *Controller) GetScheduleTypes(context *gin.Context) {
	scheduleTypes, err := c.repository.GetScheduleTypes()
	if err != nil {
		context.JSON(errors.NewError500(err, "failed get schedule types"))
		return
	}

	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Headers", "Content-Type")

	context.JSON(http.StatusOK, scheduleTypes)
}
