package filter

import "go.mongodb.org/mongo-driver/bson"

type TextSubstring struct {
	Field string `json:"field"`
	Text  string `json:"text"`
}

const TextSubstringFilterType Type = "Text"

func (n *TextSubstring) Type() Type {
	return TextSubstringFilterType
}

func NewTextSubstringFilter(
	field string,
	textSubstring string,
) *TextSubstring {
	return &TextSubstring{
		Field: field,
		Text:  textSubstring,
	}
}

func (n *TextSubstring) ToOrderedBSON() bson.D {
	return bson.D{
		{
			n.Field,
			n.Text,
		},
	}
}
