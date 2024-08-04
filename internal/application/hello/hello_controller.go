package hello

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang_fiber_base/internal/core"
	"go.uber.org/zap"
)

type HelloController interface {
	core.Router
	Hello(ctx fiber.Ctx) error
	World(ctx fiber.Ctx) error
}

type helloController struct {
	app    *fiber.App
	logger *zap.Logger
}

func NewHelloController(app *fiber.App, logger *zap.Logger) HelloController {
	return helloController{
		app:    app,
		logger: logger,
	}
}

func (h helloController) Routes() []core.Route {
	return []core.Route{
		{Method: fiber.MethodGet, Path: "/", Handler: h.Hello, Middleware: []fiber.Handler{}},
		{Method: fiber.MethodGet, Path: "/world", Handler: h.World, Middleware: []fiber.Handler{}},
	}
}

func (h helloController) Group() fiber.Router {
	return h.app.Group("/hello")
}

func (h helloController) Hello(ctx fiber.Ctx) error {
	h.logger.Info("Hello Handler")
	return ctx.SendString("World")
}

func (h helloController) World(ctx fiber.Ctx) error {
	h.logger.Info("World Handler")
	return ctx.SendString("Hello")
}
