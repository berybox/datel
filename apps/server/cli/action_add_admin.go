package main

import (
	"reflect"

	"github.com/berybox/datel/pkg/mongodb"
)

func addAdmin(URI, username, userid string) error {
	user := mongodb.User{
		Username:        username,
		ID:              userid,
		DefaultDatabase: userid,
		Collections: []mongodb.Collection{
			{
				Name:     mongodb.SettingsUsersCol,
				Database: mongodb.SettingsDB,
				Label:    "Users",
				Fields: []mongodb.Field{

					{
						Key:   getTag(&mongodb.User{}, "Username"),
						Label: "Username",
						Type: mongodb.FieldType{
							Name: "text",
						},
					},

					{
						Key:   getTag(&mongodb.User{}, "ID"),
						Label: "ID",
						Type: mongodb.FieldType{
							Name: "text",
						},
					},

					{
						Key:   getTag(&mongodb.User{}, "DefaultDatabase"),
						Label: "Default database",
						Type: mongodb.FieldType{
							Name: "text",
						},
					},
				},
			},
		},
	}

	db, err := mongodb.GetInstance(URI)
	if err != nil {
		return err
	}

	err = db.InsertOne(mongodb.SettingsDB, mongodb.SettingsUsersCol, user)
	if err != nil {
		return err
	}

	return nil
}

func getTag(v any, fieldName string) string {
	field, _ := reflect.TypeOf(v).Elem().FieldByName(fieldName)
	return field.Tag.Get("bson")
}
