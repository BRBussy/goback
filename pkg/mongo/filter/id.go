package filter

import "go.mongodb.org/mongo-driver/bson"

const IDFilterType Type = "ID"

type ID struct {
	id string
}

func NewIDFilter(id string) *ID {
	return &ID{id: id}
}

func (i *ID) Type() Type {
	return IDFilterType
}

func (i *ID) ToOrderedBSON() bson.D {
	return bson.D{
		{
			"id",
			i.id,
		},
	}
}
