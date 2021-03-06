package filter

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type TextSubstring struct {
	Field string `json:"field"`
	Text  string `json:"text"`
}

const TextSubstringFilterType Type = "TextSubstring"

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
			bson.M{
				"$regex":   fmt.Sprintf(".*%s.*", n.Text),
				"$options": "i",
			},
		},
	}
}
