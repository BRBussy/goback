package app

import (
	"github.com/BRBussy/goback/pkg/role"
	"github.com/BRBussy/goback/pkg/user"
)

var roles = []role.Role{
	{
		Name:        "User",
		Permissions: []string{},
	},
}

var rootRole = role.Role{
	Name: "Root",
	Permissions: []string{
		user.RetrieveService,
		user.ListService,
		user.AddNewUserService,
	},
}

func RoleSetup(
	roleStore role.Store,
	roleAdmin role.Admin,
) {

}
