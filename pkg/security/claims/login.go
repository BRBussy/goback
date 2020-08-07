package claims

import (
	"time"
)

const LoginClaimsType Type = "Login"

type LoginClaims struct {
	UserID         string    `validate:"required" json:"userID"`
	ExpirationTime time.Time `validate:"required" json:"expirationTime"`
}

func (l LoginClaims) Type() Type {
	return LoginClaimsType
}

func (l LoginClaims) Expired() bool {
	return time.Now().UTC().After(l.ExpirationTime)
}
