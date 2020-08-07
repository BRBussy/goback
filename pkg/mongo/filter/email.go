package filter

import "go.mongodb.org/mongo-driver/bson"

type Email struct {
	email string
}

const EmailFilterType Type = "Email"

func (e *Email) Type() Type {
	return EmailFilterType
}

func NewEmailFilter(email string) *Email {
	return &Email{email: email}
}

func (e *Email) ToOrderedBSON() bson.D {
	return bson.D{
		{
			"email",
			e.email,
		},
	}
}
