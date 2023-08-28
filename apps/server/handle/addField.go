package handle

import (
	"strings"

	"github.com/berybox/datel/apps/server/fiberutils"
	"github.com/berybox/datel/pkg/mongodb"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
)

// AddFieldGET Get page for adding a fields to the collection
func AddFieldGET(c *fiber.Ctx) error {
	_, col, _, err := fiberutils.GetUserColDB(c)
	if err != nil {
		return err
	}

	m := fiber.Map{
		"Title":      "Add field",
		"Collection": col,
		"Field":      nil,
		"Msgs":       fiberutils.CreateMessages().PullFromCtx(c),
	}
	return render("add-field", c, m)
}

// AddFieldUpdateGET Get page for updating field of the collection
func AddFieldUpdateGET(c *fiber.Ctx) error {
	_, col, _, err := fiberutils.GetUserColDB(c)
	if err != nil {
		return err
	}

	fieldIndex := col.FieldIndex(c.Params("field"))
	if fieldIndex < 0 {
		return mongodb.ErrFieldNotFound
	}

	m := fiber.Map{
		"Title":      "Edit field",
		"Collection": col,
		"Field":      col.Fields[fieldIndex],
		"Msgs":       fiberutils.CreateMessages().PullFromCtx(c),
	}
	return render("add-field", c, m)
}

// AddFieldPOST handles the HTTP POST request to add a new field.
func AddFieldPOST(c *fiber.Ctx) error {
	user, col, db, err := fiberutils.GetUserColDB(c)
	if err != nil {
		return err
	}

	f, err := c.MultipartForm()
	if err != nil {
		return err
	}

	newField := multipartToField(f.Value)
	fieldIndex := col.FieldIndex(newField.Key)
	if fieldIndex >= 0 {
		col.Fields[fieldIndex] = newField

		fiberutils.CreateMessages().AddSuccess("Field has been updated").PutToCtx(c)
	} else {
		col.Fields = append(col.Fields, newField)

		fiberutils.CreateMessages().AddSuccess("Field has been added").PutToCtx(c)
	}

	err = db.UserUpdateCollection(&user, col)
	if err != nil {
		return err
	}

	return AddFieldGET(c)
}

func multipartToField(req map[string][]string) mongodb.Field {
	ret := mongodb.Field{
		Key:   strings.Join(req["field-id"], ""),
		Label: strings.Join(req["field-label"], ""),
		Type: mongodb.FieldType{
			Name: strings.Join(req["field-type"], ""),
		},
	}

	omitKeys := []string{
		"field-id",
		"field-label",
		"field-type",
	}

	props := make(map[string]interface{})
	for key, val := range req {
		if slices.Contains(omitKeys, key) {
			continue
		}
		key = strings.TrimPrefix(key, "field-")
		props[key] = strings.Join(val, "")
	}

	ret.Type.Properties = props

	return ret
}
