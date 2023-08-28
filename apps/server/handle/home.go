package handle

import (
	"github.com/berybox/datel/apps/server/fiberutils"
	"github.com/gofiber/fiber/v2"
)

// HomeGET Get homepage with basic overview of user
func HomeGET(c *fiber.Ctx) error {
	user, _, err := fiberutils.GetUserDB(c)
	if err != nil {
		return err
	}

	msgs := fiberutils.CreateMessages().PullFromCtx(c)

	m := fiber.Map{
		"User": user,
		"Msgs": msgs,
	}
	return render("home", c, m)
}
