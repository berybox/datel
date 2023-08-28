package handle

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// ErrorGET handles custom error page
func ErrorGET(c *fiber.Ctx, err error) error {
	if isHxReqest(c) {
		c.Set("HX-Retarget", "body")
		c.Set("HX-Reswap", "outerHTML")
		c.Set("HX-Push-Url", "/error")
	}

	m := fiber.Map{
		"Text": fmt.Sprintf("%s", err),
	}

	err = c.Render("error", m, "error")
	if err != nil {
		return err
	}

	return nil
}
