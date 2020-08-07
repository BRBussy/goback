package filter

import "go.mongodb.org/mongo-driver/bson"

type Username struct {
	username string
}

const UsernameFilterType Type = "Username"

func (e *Username) Type() Type {
	return UsernameFilterType
}

func NewUsernameFilter(username string) *Username {
	return &Username{username: username}
}

func (e *Username) ToOrderedBSON() bson.D {
	return bson.D{
		{
			"username",
			e.username,
		},
	}
}
