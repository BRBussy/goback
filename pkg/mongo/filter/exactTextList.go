package filter

import "go.mongodb.org/mongo-driver/bson"

const ExactTextListFilterType Type = "ExactTextList"

type ExactTextListFilter struct {
	field string
	list  []string
}

func NewExactTextListFilter(
	field string,
	list []string,
) *ExactTextListFilter {
	return &ExactTextListFilter{
		field: field,
		list:  list,
	}
}

func (e *ExactTextListFilter) Type() Type {
	return ExactTextListFilterType
}

func (e *ExactTextListFilter) ToOrderedBSON() bson.D {
	return bson.D{
		{
			e.field,
			bson.M{"$in": e.list},
		},
	}
}
