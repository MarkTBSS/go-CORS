package servers

import (
	"github.com/MarkTBSS/go-CORS/modules/middlewares/middlewaresHandlers"
	"github.com/MarkTBSS/go-CORS/modules/middlewares/middlewaresRepositories"
	"github.com/MarkTBSS/go-CORS/modules/middlewares/middlewaresUsecases"
	"github.com/MarkTBSS/go-CORS/modules/monitor/monitorHandlers"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	r fiber.Router
	s *server
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r: r,
		s: s,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.NewMiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.NewMiddlewaresUsecase(repository)
	return middlewaresHandlers.NewMiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.NewMonitorHandler(m.s.cfg)
	m.r.Get("/", handler.HealthCheck)
}
