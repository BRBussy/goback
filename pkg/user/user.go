package user

type User struct {
	ID       string   `json:"id" bson:"id"`
	Email    string   `validate:"required,email" json:"email" bson:"email"`
	RoleIDs  []string `validate:"required" json:"roleIDs" bson:"roleIDs"`
	Password []byte   `json:"-" bson:"password"`
}

func (u User) Equal(u2 User) bool {
	if u.ID != u2.ID {
		return false
	}
	if u.Email != u2.ID {
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

	// map to store occurrences in u.Password
	uPasswordByteCount := make(map[byte]int)
	for _, b := range u.Password {
		uPasswordByteCount[b]++
	}

	// check u2
	for _, b := range u2.Password {
		count, found := uPasswordByteCount[b]
		if !found || count == 0 {
			return false
		}
		uPasswordByteCount[b]--
	}

	return true
}
