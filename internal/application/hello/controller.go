package hello

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang_fiber_base/internal/core"
	"go.uber.org/zap"
)

type Controller interface {
	core.Router
	Hello(ctx *fiber.Ctx) error
	World(ctx *fiber.Ctx) error
}

type helloController struct {
	app     *fiber.App
	logger  *zap.Logger
	service Service
}

func NewHelloController(p core.Param, service Service) Controller {
	return helloController{
		app:     p.App,
		logger:  p.Logger,
		service: service,
	}
}

func (h helloController) Routes() []core.Route {
	return []core.Route{
		{Method: fiber.MethodGet, Path: "/", Handler: h.Hello},
		{Method: fiber.MethodGet, Path: "/world", Handler: h.World},
	}
}

func (h helloController) Group() fiber.Router {
	return h.app.Group("/hello")
}

func (h helloController) Hello(ctx *fiber.Ctx) error {
	return ctx.SendString(h.service.Hello())
}

func (h helloController) World(ctx *fiber.Ctx) error {
	return ctx.SendString(h.service.World())
}
