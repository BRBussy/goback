package filter

import "go.mongodb.org/mongo-driver/bson"

type Name struct {
	name string
}

const NameFilterType Type = "Name"

func (n *Name) Type() Type {
	return NameFilterType
}

func NewNameFilter(name string) *Name {
	return &Name{name: name}
}

func (n *Name) ToOrderedBSON() bson.D {
	return bson.D{
		{
			"name",
			n.name,
		},
	}
}
