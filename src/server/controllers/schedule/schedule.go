package schedule

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rsreu-class-schedule/server/errors"
)

type GetRequest struct {
	Type string `json:"type"`
}

func (r *GetRequest) ValidateRequest() error {
	if r.Type == "" {
		return fmt.Errorf("empty schedule type")
	}
	return nil
}

func (c *Controller) GetSchedule(context *gin.Context) {
	req := new(GetRequest)

	if err := c.unmarshalRequest(context.Request, req); err != nil {
		context.IndentedJSON(errors.NewError400(err, "failed get request"))
		return
	}

	context.IndentedJSON(http.StatusOK, req)
}
