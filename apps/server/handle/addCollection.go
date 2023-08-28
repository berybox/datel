package handle

import (
	"strings"

	"github.com/berybox/datel/apps/server/fiberutils"
	"github.com/berybox/datel/pkg/mongodb"
	"github.com/gofiber/fiber/v2"
)

// AddCollectionGET Get page for adding a collection
func AddCollectionGET(c *fiber.Ctx) error {
	user, _, err := fiberutils.GetUserDB(c)
	if err != nil {
		return err
	}

	m := fiber.Map{
		"Title": "Add collection",
		"User":  user,
	}
	return render("add-collection", c, m)
}

// AddCollectionPOST Add collection to the user settings
func AddCollectionPOST(c *fiber.Ctx) error {
	user, db, err := fiberutils.GetUserDB(c)
	if err != nil {
		return err
	}

	f, err := c.MultipartForm()
	if err != nil {
		return err
	}

	newCol := mongodb.Collection{
		Name:     strings.Join(f.Value["collection-name"], ""),
		Label:    strings.Join(f.Value["collection-label"], " "),
		Database: strings.Join(f.Value["collection-database"], ""),
	}

	msgs := fiberutils.CreateMessages()

	colExists, err := db.HasCollection(newCol.Database, newCol.Name)
	if err != nil {
		return err
	}

	if colExists || user.CollectionIndex(newCol.Name) >= 0 {
		msgs.AddDanger("Collection already exists")
	} else {
		err := db.UserAddCollection(&user, newCol)
		if err != nil {
			return err
		}

		msgs.AddSuccess("Collection has been added")
	}

	m := fiber.Map{
		"Title": "Add collection",
		"User":  user,
		"Msgs":  msgs,
	}
	return render("add-collection", c, m)
}

// DeleteCollectionGET Delete collection
func DeleteCollectionGET(c *fiber.Ctx) error {
	user, col, db, err := fiberutils.GetUserColDB(c)
	if err != nil {
		return err
	}

	err = db.UserDeleteCollection(&user, col)
	if err != nil {
		return err
	}

	err = db.DropCollection(col.Database, col.Name)
	if err != nil {
		return err
	}

	fiberutils.CreateMessages().AddSuccess("Collection has been deleted").PutToCtx(c)
	setState(c, state{URL: "/"})

	return HomeGET(c)
}
