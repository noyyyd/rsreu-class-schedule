package server

import (
	"github.com/gin-gonic/gin"
	"rsreu-class-schedule/config"
	"rsreu-class-schedule/server/controllers/schedule"
)

type Server struct {
	router *gin.Engine

	config *config.Config

	schedule *schedule.Controller
}

func New(config *config.Config) *Server {
	return &Server{
		config: config,

		router: gin.Default(),

		schedule: schedule.NewController(config),
	}
}

func (s *Server) Run() {
	s.initRoutes()

	s.router.Run()
}

func (s *Server) initRoutes() {
	s.router.GET("/schedule_types", s.schedule.GetScheduleTypes)
	s.router.POST("/schedule", s.schedule.GetSchedule)
}
