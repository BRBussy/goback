package filter

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

type Filter interface {
	ToOrderedBSON() bson.D
	Type() Type
}

type Type string

func (t Type) String() string {
	return string(t)
}

type SerializedFilter struct {
	Filter Filter
}

type typeHolder struct {
	Type Type `json:"type"`
}

func (s SerializedFilter) UnmarshalJSON(bytes []byte) error {
	var t typeHolder
	if err := json.Unmarshal(bytes, &t); err != nil {
		return NewErrJSONUnmarshallError(err)
	}

	switch t.Type {
	case EmailFilterType:
		f := new(Email)
		if err := json.Unmarshal(bytes, f); err != nil {
			return NewErrJSONUnmarshallError(err)
		}
		s.Filter = f

	case IDFilterType:
		f := new(ID)
		if err := json.Unmarshal(bytes, f); err != nil {
			return NewErrJSONUnmarshallError(err)
		}
		s.Filter = f

	default:
		return NewErrInvalidType(t.Type)
	}

	return nil
}
