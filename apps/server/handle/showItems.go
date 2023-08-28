package handle

import (
	"github.com/berybox/datel/apps/server/fiberutils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// ShowItemsGET Get page with items inside collection
func ShowItemsGET(c *fiber.Ctx) error {
	_, col, db, err := fiberutils.GetUserColDB(c)
	if err != nil {
		return err
	}

	items, err := db.GetItems(col, bson.D{})
	if err != nil {
		return err
	}

	m := fiber.Map{
		"Title": "Items list",
		"Items": items,
	}
	return render("show-items", c, m)
}
