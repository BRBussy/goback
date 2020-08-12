package claims

import (
	"time"
)

type Claims struct {
	UserID         string    `validate:"required" json:"userID"`
	ExpirationTime time.Time `validate:"required" json:"expirationTime"`
}

func NewClaims(
	userID string,
	expirationTime time.Time,
) *Claims {
	return &Claims{UserID: userID, ExpirationTime: expirationTime}
}

func (c *Claims) Expired() bool {
	return time.Now().UTC().After(c.ExpirationTime)
}
