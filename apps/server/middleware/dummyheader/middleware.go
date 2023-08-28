package dummyheader

import (
	"github.com/gofiber/fiber/v2"
)

// Config necessary for Fiber
type Config struct {
	Key   string
	Value string
}

// New Fiber middleware, header(s) from config Key/Value to the request
func New(config ...Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		for _, conf := range config {
			c.Request().Header.Add(conf.Key, conf.Value)
		}
		return c.Next()
	}
}
