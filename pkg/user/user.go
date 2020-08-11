package user

type User struct {
	ID         string   `json:"id" bson:"id"`
	Username   string   `json:"username" bson:"username"`
	Email      string   `validate:"required,email" json:"email" bson:"email"`
	RoleIDs    []string `validate:"required" json:"roleIDs" bson:"roleIDs"`
	Registered bool     `json:"registered" bson:"registered"`
	Password   []byte   `json:"-" bson:"password"`
}

// Equal compares 2 users - excluding passwords
func (u User) Equal(u2 User) bool {
	if u.ID != u2.ID {
		return false
	}
	if u.Username != u2.Username {
		return false
	}
	if u.Email != u2.Email {
		return false
	}
	if u.Registered != u2.Registered {
		return false
	}

	if len(u.RoleIDs) != len(u2.RoleIDs) {
		return false
	}

	// map to store occurrences in u.RoleIDs
	uRoleIDCount := make(map[string]int)
	for _, roleID := range u.RoleIDs {
		uRoleIDCount[roleID]++
	}

	// check u2
	for _, roleID := range u2.RoleIDs {
		count, found := uRoleIDCount[roleID]
		if !found || count == 0 {
			return false
		}
		uRoleIDCount[roleID]--
	}

	return true
}
