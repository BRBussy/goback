package filter

import "go.mongodb.org/mongo-driver/bson"

type TextExact struct {
	Field     string `json:"field"`
	TextExact string `json:"textExact"`
}

const TextExactFilterType Type = "TextExact"

func (n *TextExact) Type() Type {
	return TextExactFilterType
}

func NewTextExactFilter(field, textExact string) *TextExact {
	return &TextExact{
		Field:     field,
		TextExact: textExact,
	}
}

func (n *TextExact) ToOrderedBSON() bson.D {
	return bson.D{
		{
			n.Field,
			n.TextExact,
		},
	}
}
