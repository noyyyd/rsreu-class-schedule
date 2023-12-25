package schedule

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"rsreu-class-schedule/config"
	"rsreu-class-schedule/repository"
	"rsreu-class-schedule/server/controllers"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	config *config.Config

	repository repository.ScheduleRepository
}

func NewController(config *config.Config) *Controller {
	return &Controller{
		config: config,

		repository: repository.NewRepository(config),
	}
}

func (c *Controller) unmarshalRequest(context *gin.Context, dst interface{}) error {
	if dstType := reflect.TypeOf(dst); dstType.Kind() != reflect.Ptr && dstType.Kind() != reflect.Map {
		return fmt.Errorf("dst struct is not pointer or map is %T", dst)
	}

	body, err := io.ReadAll(context.Request.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	err = json.Unmarshal(body, dst)
	if err != nil {
		log.Println(err)
		return err
	}

	if validator, ok := dst.(controllers.RequestValidator); ok {
		err = validator.ValidateRequest()
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
