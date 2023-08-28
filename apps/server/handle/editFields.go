package handle

import (
	"github.com/berybox/datel/apps/server/fiberutils"
	"github.com/berybox/datel/pkg/mongodb"
	"github.com/gofiber/fiber/v2"
	"github.com/tidwall/gjson"
)

// EditFieldsGET Get page for editing of fields
func EditFieldsGET(c *fiber.Ctx) error {
	_, col, _, err := fiberutils.GetUserColDB(c)
	if err != nil {
		return err
	}

	m := fiber.Map{
		"Title":      "Edit fields",
		"Collection": col,
		"Msgs":       fiberutils.CreateMessages().PullFromCtx(c),
	}
	return render("edit-fields", c, m)
}

// EditFieldsPOST Edit fields
func EditFieldsPOST(c *fiber.Ctx) error {
	user, col, db, err := fiberutils.GetUserColDB(c)
	if err != nil {
		return err
	}

	newCol := mongodb.Collection{
		Name:     col.Name,
		Label:    col.Label,
		Database: col.Database,
	}

	req := gjson.ParseBytes(c.Body())
	req.ForEach(func(key, value gjson.Result) bool {
		i := col.FieldIndex(key.String())
		if i >= 0 {
			newCol.Fields = append(newCol.Fields, col.Fields[i])
		}
		return true
	})

	err = db.UserUpdateCollection(&user, newCol)
	if err != nil {
		return err
	}

	fiberutils.CreateMessages().AddSuccess("Fields editing saved").PutToCtx(c)

	return EditFieldsGET(c)
}
