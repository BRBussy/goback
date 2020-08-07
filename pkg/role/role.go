package role

type Role struct {
	ID          string   `validate:"required" json:"id" bson:"id"`
	Name        string   `validate:"required" json:"name" bson:"name"`
	Permissions []string `validate:"required" json:"permissions" bson:"permissions"`
}

func (r *Role) AddUniquePermissions(permissionsToAdd ...string) {
	// index existing permissions
	existingPermIdx := make(map[string]bool)
	for _, existingPerm := range r.Permissions {
		existingPermIdx[existingPerm] = true
	}

	// add permissions not already on role
	for _, permToAdd := range permissionsToAdd {
		if !existingPermIdx[permToAdd] {
			r.Permissions = append(
				r.Permissions,
				permToAdd,
			)
		}
	}
}

func (r *Role) Equal(r2 Role) bool {
	if r.ID != r2.ID {
		return false
	}
	if r.Name != r2.Name {
		return false
	}
	if len(r.Permissions) != len(r2.Permissions) {
		return false
	}

	// map to store occurrences in r.Permissions
	rPermHashCount := make(map[string]int)
	for _, perm := range r.Permissions {
		rPermHashCount[perm]++
	}

	// check r2
	for _, perm := range r2.Permissions {
		count, found := rPermHashCount[perm]
		if !found || count == 0 {
			return false
		}
		rPermHashCount[perm]--
	}

	return true
}
