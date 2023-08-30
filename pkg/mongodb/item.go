package mongodb

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// LabelTypeName Name of the key for the "label", i.e. gets priority for display
	LabelTypeName = "_label"
)

var (
	// ErrItemNotFound item not found error
	ErrItemNotFound = errors.New("item not found")
)

// Item item (mongo document) in collection
type Item struct {
	ID         string      `json:"_id" bson:"_id"`
	Collection *Collection `json:"-" bson:"-"`
	Label      string      `json:"-" bson:"-"`
	Data       []ItemData  `json:"data" bson:"data"`
}

// ItemData items of the document
type ItemData struct {
	Field Field       `json:"field" bson:"field"`
	Value interface{} `json:"value" bson:"value"`
}

// EmptyItem creates an empty item based on the collection fields
func EmptyItem(col Collection) Item {
	ret := Item{
		ID:         "",
		Collection: &col,
		Label:      "",
	}
	for _, field := range col.Fields {
		ret.Data = append(ret.Data, ItemData{Field: field})
	}
	return ret
}

// ItemFromMultipart converts a multipart request to an item
func ItemFromMultipart(col Collection, multipart *multipart.Form) Item {
	ret := Item{
		Collection: &col,
	}

	if val, ok := multipart.Value["_id"]; ok {
		if len(val) > 0 {
			ret.ID = val[0]
		}
	}

	if val, ok := multipart.Value[LabelTypeName]; ok {
		ret.Label = strings.Join(val, " ")
	}

	for _, field := range col.Fields {
		val := multipartValueToType(multipart.Value[field.Key], field.Type)
		ret.Data = append(ret.Data, ItemData{
			Field: field,
			Value: val,
		})
	}

	return ret
}

func multipartValueToType(val []string, typ FieldType) interface{} {
	var ret interface{}

	switch typ.Name {
	case "text":
		//ret = strings.Join(val, " ")
		if len(val) > 0 {
			ret = val[0]
		}
	case "number":
		ret, _ = strconv.Atoi(strings.Join(val, ""))
	default:
		ret = val
	}

	return ret
}

func (t *Item) toBSON() bson.D {
	var ret bson.D

	/*
		if t.ID != "" {
			ret = append(ret, primitive.E{Key: "_id", Value: t.ID})
		}
	*/

	for _, data := range t.Data {
		ret = append(ret, primitive.E{Key: data.Field.Key, Value: data.Value})
	}

	return ret
}

// AddItem adds an entry to the database, if ID is "", database will created it
func (m *Instance) AddItem(item Item) error {
	err := m.InsertOne(item.Collection.Database, item.Collection.Name, item.toBSON())
	if err != nil {
		return err
	}
	return nil
}

// DeleteItem deletes item which gets filtered by the item.ID
func (m *Instance) DeleteItem(item Item) error {
	filter := IDFilter(item.ID)

	_, err := m.DeleteOne(item.Collection.Database, item.Collection.Name, filter)
	if err != nil {
		return err
	}

	return nil
}

// DeleteItemByID deletes item by ID
func (m *Instance) DeleteItemByID(col Collection, id string) error {
	item := Item{
		ID:         id,
		Collection: &col,
	}

	err := m.DeleteItem(item)
	if err != nil {
		return err
	}

	return nil
}

// ReplaceItem replaces the item which is filtered by ID
func (m *Instance) ReplaceItem(item Item) error {
	filter := IDFilter(item.ID)

	_, err := m.ReplaceOne(item.Collection.Database, item.Collection.Name, filter, item.toBSON())
	if err != nil {
		return err
	}

	return nil
}

// UpdateItem updates item in the collection. Existing fields are kept.
func (m *Instance) UpdateItem(item Item) error {
	filter := IDFilter(item.ID)

	updateData, err := bson.Marshal(item)
	if err != nil {
		return err
	}

	var update bson.M
	err = bson.Unmarshal(updateData, &update)
	if err != nil {
		return err
	}

	delete(update, "_id")

	_, err = m.UpdateOne(item.Collection.Database, item.Collection.Name, filter, bson.M{"$set": item.toBSON()})
	if err != nil {
		return err
	}

	return nil
}

// GetItems gets items from the database
func (m *Instance) GetItems(col Collection, filter interface{}) ([]Item, error) {
	var rets []map[string]interface{}
	err := m.UnmarshallCollection(col.Database, col.Name, filter, &rets)
	if err != nil {
		return []Item{}, err
	}

	if len(rets) < 1 {
		return []Item{}, ErrItemNotFound
	}

	items := itemsFromMongoDoc(col, rets)

	return items, nil
}

func itemsFromMongoDoc(col Collection, docs []map[string]interface{}) []Item {
	var ret []Item

	for _, doc := range docs {
		var id string
		if oid, ok := doc["_id"].(primitive.ObjectID); ok {
			id = oid.Hex()
		} else {
			id = doc["_id"].(string)
		}

		item := Item{
			Collection: &col,
			ID:         id,
		}

		for _, f := range col.Fields {
			item.Data = append(item.Data, ItemData{Field: f, Value: doc[f.Key]})
		}

		if doc[LabelTypeName] != nil {
			item.Label = fmt.Sprintf("%s", doc[LabelTypeName])
		}
		if item.Label == "" && len(item.Data) > 0 {
			item.Label = fmt.Sprintf("%s", item.Data[0].Value)
		}

		ret = append(ret, item)
	}

	return ret
}
