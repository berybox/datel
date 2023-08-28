package handle

import (
	"github.com/berybox/datel/apps/server/fiberutils"
	"github.com/gofiber/fiber/v2"
)

var (
	defaultLayouts = []string{"layout"}
)

type state struct {
	URL string
}

func render(page string, c *fiber.Ctx, m fiber.Map) error {
	var layouts []string
	if !isHxReqest(c) {
		layouts = append(layouts, defaultLayouts...)
	}

	//if c.Method() == fiber.MethodGet {
	c.Set("HX-Push-Url", c.OriginalURL())
	//}

	menu, err := fiberutils.GetMenu(c)
	if err != nil {
		return err
	}

	m["Menu"] = menu
	m["Username"] = fiberutils.GetUserID(c)
	m["HTMLState"] = getState(c)

	err = c.Render(page, m, layouts...)
	if err != nil {
		return err
	}

	return nil
}

func isHxReqest(c *fiber.Ctx) bool {
	ret := c.Get("Hx-Request")
	return ret == "true"
}

func setState(c *fiber.Ctx, state state) {
	c.Locals("HTMLState", state)
}

func getState(c *fiber.Ctx) *state {
	if s, ok := c.Locals("HTMLState").(state); ok {
		return &s
	}
	return nil
}
