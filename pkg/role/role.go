package role

type Role struct {
	ID          string   `validate:"required" json:"id" bson:"id"`
	Name        string   `validate:"required" json:"name" bson:"name"`
	Permissions []string `validate:"required" json:"permissions" bson:"permissions"`
}
