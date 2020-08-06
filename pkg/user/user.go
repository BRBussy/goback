package user

type User struct {
	ID       string   `json:"id" bson:"id"`
	Email    string   `validate:"required,email" json:"email" bson:"email"`
	RoleIDs  []string `validate:"required" json:"roleIDs" bson:"roleIDs"`
	Password []byte   `json:"-" bson:"password"`
}
