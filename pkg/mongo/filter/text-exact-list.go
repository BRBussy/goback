package filter

import "go.mongodb.org/mongo-driver/bson"

const TextExactListFilterType Type = "TextExactList"

type TextExactListFilter struct {
	Field string   `json:"Field"`
	List  []string `json:"list"`
}

func NewTextExactListFilter(
	field string,
	list []string,
) *TextExactListFilter {
	return &TextExactListFilter{
		Field: field,
		List:  list,
	}
}

func (e *TextExactListFilter) Type() Type {
	return TextExactListFilterType
}

func (e *TextExactListFilter) ToOrderedBSON() bson.D {
	return bson.D{
		{
			e.Field,
			bson.M{"$in": e.List},
		},
	}
}
