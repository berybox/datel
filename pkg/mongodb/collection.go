package mongodb

import (
	"errors"
)

var (
	//ErrCollectionNotFound collection not found error
	ErrCollectionNotFound = errors.New("collection not found")

	//ErrFieldNotFound field not found error
	ErrFieldNotFound = errors.New("field not found")
)

// Collection collection in database and its fields
type Collection struct {
	Name     string  `json:"name" bson:"name"`
	Label    string  `json:"label" bson:"label"`
	Database string  `json:"database" bson:"database"`
	Fields   []Field `json:"fields" bson:"fields"`
}

// Field fields of collection
type Field struct {
	Key   string    `json:"key" bson:"key"`
	Label string    `json:"label" bson:"label"`
	Type  FieldType `json:"type" bson:"type"`
}

// FieldType type of field
type FieldType struct {
	Name       string                 `json:"name" bson:"name"`
	Properties map[string]interface{} `json:"properties,omitempty" bson:"properties,omitempty"`
}

// FieldIndex find field index by key
func (c *Collection) FieldIndex(key string) int {
	for i, f := range c.Fields {
		if f.Key == key {
			return i
		}
	}
	return -1
}
