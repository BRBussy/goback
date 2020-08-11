package filter

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

type Type string

func (t Type) String() string {
	return string(t)
}

type Filter interface {
	ToOrderedBSON() bson.D
	Type() Type
}

type SerializedFilter struct {
	Filter Filter
}

type typeHolder struct {
	Type Type `json:"type"`
}

func (s *SerializedFilter) UnmarshalJSON(bytes []byte) error {
	var t typeHolder
	if err := json.Unmarshal(bytes, &t); err != nil {
		return NewErrJSONUnmarshallError(err)
	}

	switch t.Type {
	case TextExactFilterType:
		f := new(TextExact)
		if err := json.Unmarshal(bytes, f); err != nil {
			return NewErrJSONUnmarshallError(err)
		}
		s.Filter = f

	case TextExactListFilterType:
		f := new(TextExactListFilter)
		if err := json.Unmarshal(bytes, f); err != nil {
			return NewErrJSONUnmarshallError(err)
		}
		s.Filter = f

	case TextSubstringFilterType:
		f := new(TextSubstring)
		if err := json.Unmarshal(bytes, f); err != nil {
			return NewErrJSONUnmarshallError(err)
		}
		s.Filter = f

	default:
		return NewErrInvalidType(t.Type)
	}

	return nil
}
