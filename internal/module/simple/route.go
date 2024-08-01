package simple

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang_fiber_base/internal/app"
)

func Get(c fiber.Ctx) error {
	app_ := app.FromLocals(c)
	app_.Logger.Info("Hello Simple")
	return c.SendString("Hello Simple")
}

func RegisterRoute(router fiber.Router) {
	router.Get("/", Get)
}
