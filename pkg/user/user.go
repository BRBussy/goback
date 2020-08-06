package user

type User struct {
	ID       string `json:"id" bson:"id"`
	Name     string `validate:"required" json:"name" bson:"name"`
	Email    string `validate:"required,email" json:"email" bson:"email"`
	RoleIDs  string `validate:"required" json:"roleIDs" bson:"roleIDs"`
	Password []byte `json:"-" bson:"password"`
}
