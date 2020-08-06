package claims

import (
	"encoding/json"
	"time"
)

type Login struct {
	UserID         string `validate:"required" json:"userID"`
	ExpirationTime int64  `validate:"required" json:"expirationTime"`
}

func (l Login) Type() Type {
	return LoginClaimsType
}

func (l Login) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type           Type   `json:"type"`
		UserID         string `json:"userID"`
		ExpirationTime int64  `json:"expirationTime"`
	}{
		Type:           l.Type(),
		UserID:         l.UserID,
		ExpirationTime: l.ExpirationTime,
	})
}

func (l Login) Expired() bool {
	return time.Now().UTC().After(time.Unix(l.ExpirationTime, 0).UTC())
}
