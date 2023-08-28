package useroverride

import (
	"github.com/gofiber/fiber/v2"
)

var (
	// FieldName name of Fiber Ctx variable field
	FieldName = "UsernameOverride"
)

// Config necessary for Fiber
type Config struct {
	Name string
}

// New Fiber middleware, puts username into the local varibles of Fiber Ctx
func New(config ...Config) fiber.Handler {
	if len(config) > 0 {
		return func(c *fiber.Ctx) error {
			c.Locals(FieldName, config[0].Name)
			return c.Next()
		}
	}
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
