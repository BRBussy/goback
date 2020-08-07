package filter

import "go.mongodb.org/mongo-driver/bson"

type EmailAddress struct {
	emailAddress string
}

const EmailAddressFilterType Type = "EmailAddress"

func (e *EmailAddress) Type() Type {
	return EmailAddressFilterType
}

func NewEmailAddressFilter(emailAddress string) *EmailAddress {
	return &EmailAddress{emailAddress: emailAddress}
}

func (e *EmailAddress) ToOrderedBSON() bson.D {
	return bson.D{
		{
			"emailAddress",
			e.emailAddress,
		},
	}
}
