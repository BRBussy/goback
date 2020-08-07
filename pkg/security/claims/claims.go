package claims

import "encoding/json"

type Type string

func (t Type) String() string {
	return string(t)
}

type Claims interface {
	Type() Type    // Returns the Type of the claims
	Expired() bool // Returns whether or not the claims are expired
}

type SerializedClaims struct {
	Claims Claims `json:"claims"`
	Type   Type   `json:"type"`
}

type unmarshalHolder struct {
	Type       Type            `json:"type"`
	ClaimsData json.RawMessage `json:"claims"`
}

func (s *SerializedClaims) MarshalJSON() ([]byte, error) {
	s.Type = s.Claims.Type()
	return json.Marshal(s)
}

func (s *SerializedClaims) UnmarshalJSON(bytes []byte) error {
	var h unmarshalHolder
	if err := json.Unmarshal(bytes, &h); err != nil {
		return NewErrJSONUnmarshallError(err)
	}

	switch h.Type {
	case LoginClaimsType:
		c := new(LoginClaims)
		if err := json.Unmarshal(h.ClaimsData, c); err != nil {
			return NewErrJSONUnmarshallError(err)
		}
		s.Claims = c

	default:
		return NewErrInvalidType(h.Type)
	}

	return nil
}
