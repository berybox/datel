package fiberutils

import (
	"github.com/berybox/datel/pkg/mongodb"
	"github.com/gofiber/fiber/v2"
)

// MenuGroup group of menu items
type MenuGroup struct {
	Label     string
	DeleteURL string
	Items     []MenuItem
}

// MenuItem single item in menu
type MenuItem struct {
	Label string
	URL   string
}

// GetMenu gets menu items for main layout
func GetMenu(c *fiber.Ctx) ([]MenuGroup, error) {
	var ret []MenuGroup

	user, _, err := GetUserDB(c)
	if err != nil {
		return ret, err
	}

	for _, col := range user.Collections {
		m := MenuGroup{
			Label: col.Label,
			Items: []MenuItem{
				{
					Label: "Show items",
					URL:   "/show-items/" + col.Name,
				},
				{
					Label: "Add item",
					URL:   "/add-item/" + col.Name,
				},
				{
					Label: "Add field",
					URL:   "/add-field/" + col.Name,
				},
				{
					Label: "Edit fields",
					URL:   "/edit-fields/" + col.Name,
				},
			},
		}

		if col.Database != mongodb.SettingsDB {
			m.DeleteURL = "/delete-collection/" + col.Name
		}

		ret = append(ret, m)
	}

	return ret, nil
}
