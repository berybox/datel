package fiberutils

import (
	"github.com/berybox/datel/apps/server/middleware/useroverride"
	"github.com/berybox/datel/pkg/mongodb"
	"github.com/gofiber/fiber/v2"
)

const (
	// UserIDHeader name of header which passes user ID in request
	UserIDHeader = "X-UserID"
)

// GetUserID gets ID of current user
func GetUserID(c *fiber.Ctx) string {
	var userid string
	u := c.Locals(useroverride.FieldName)
	if u == nil {
		userid = c.Get(UserIDHeader)
	} else {
		var ok bool
		userid, ok = u.(string)
		if !ok {
			userid = ""
		}
	}
	return userid
}

// GetUserColDB shortcut, frequently used
func GetUserColDB(c *fiber.Ctx) (mongodb.User, mongodb.Collection, *mongodb.Instance, error) {
	user, db, err := GetUserDB(c)
	if err != nil {
		return mongodb.User{}, mongodb.Collection{}, &mongodb.Instance{}, err
	}

	col, err := user.CollectionByName(c.Params("collection"))
	if err != nil {
		return mongodb.User{}, mongodb.Collection{}, &mongodb.Instance{}, err
	}

	return user, col, db, nil
}

// GetUserDB shortcut, frequently used
func GetUserDB(c *fiber.Ctx) (mongodb.User, *mongodb.Instance, error) {
	db, err := mongodb.GetInstance(mongodb.CurrentURI())
	if err != nil {
		return mongodb.User{}, &mongodb.Instance{}, err
	}

	user, err := db.GetUser(GetUserID(c))
	if err != nil {
		return mongodb.User{}, &mongodb.Instance{}, err
	}

	return user, db, nil
}
