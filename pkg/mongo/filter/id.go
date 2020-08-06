package filter

import "go.mongodb.org/mongo-driver/bson"

type ID struct {
	id string
}

const IDFilterType Type = "ID"

func (i *ID) Type() Type {
	return IDFilterType
}

func NewIDFilter(id string) *ID {
	return &ID{id: id}
}

func (i *ID) ToOrderedBSON() bson.D {
	return bson.D{
		{
			"id",
			i.id,
		},
	}
}
