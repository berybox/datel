package handle

import (
	"github.com/berybox/datel/apps/server/fiberutils"
	"github.com/berybox/datel/pkg/mongodb"
	"github.com/gofiber/fiber/v2"
)

// AddItemGET Get page for adding items
func AddItemGET(c *fiber.Ctx) error {
	_, col, _, err := fiberutils.GetUserColDB(c)
	if err != nil {
		return err
	}

	m := fiber.Map{
		"Title": "Add item",
		"Item":  mongodb.EmptyItem(col),
		"Msgs":  fiberutils.CreateMessages().PullFromCtx(c),
	}
	return render("add-item", c, m)
}

// AddItemUpdateGET Get page for updating an item
func AddItemUpdateGET(c *fiber.Ctx) error {
	_, col, db, err := fiberutils.GetUserColDB(c)
	if err != nil {
		return err
	}

	filter := mongodb.IDFilter(c.Params("dbid"))

	items, err := db.GetItems(col, filter)
	if err != nil {
		return err
	}

	if len(items) < 1 {
		return mongodb.ErrItemNotFound
	}

	m := fiber.Map{
		"Title": "Edit item",
		"Item":  items[0],
		"Msgs":  fiberutils.CreateMessages().PullFromCtx(c),
	}
	return render("add-item", c, m)
}

// AddItemPOST Add item to the collection
func AddItemPOST(c *fiber.Ctx) error {
	_, col, db, err := fiberutils.GetUserColDB(c)
	if err != nil {
		return err
	}

	f, err := c.MultipartForm()
	if err != nil {
		return err
	}

	item := mongodb.ItemFromMultipart(col, f)

	err = db.AddItem(item)
	if err != nil {
		return err
	}

	fiberutils.CreateMessages().AddSuccess("Item has been added").PutToCtx(c)
	return AddItemGET(c)
}

// AddItemUpdatePOST Update item in the collection
func AddItemUpdatePOST(c *fiber.Ctx) error {
	_, col, db, err := fiberutils.GetUserColDB(c)
	if err != nil {
		return err
	}

	f, err := c.MultipartForm()
	if err != nil {
		return err
	}

	item := mongodb.ItemFromMultipart(col, f)

	err = db.ReplaceItem(item)
	if err != nil {
		return err
	}

	fiberutils.CreateMessages().AddSuccess("Item has been updated").PutToCtx(c)
	return AddItemUpdateGET(c)
}

// AddItemDELETE item from collection
func AddItemDELETE(c *fiber.Ctx) error {
	_, col, db, err := fiberutils.GetUserColDB(c)
	if err != nil {
		return err
	}

	err = db.DeleteItemByID(col, c.Params("dbid"))
	if err != nil {
		return err
	}

	return ShowItemsGET(c)
}
