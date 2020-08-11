package setup

import (
	"errors"
	"fmt"
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/mongo/filter"
	"github.com/BRBussy/goback/pkg/role"
	"github.com/BRBussy/goback/pkg/user"
	"github.com/rs/zerolog/log"
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

func RoleSync(
	roleStore role.Store,
	roleAdmin role.Admin,
) {
	// add all permissions to the root role
	for _, r := range roles {
		rootRole.AddUniquePermissions(r.Permissions...)
	}

	// add root role to all roles
	roles = append(
		roles,
		rootRole,
	)

	// synchronise roles
	for _, roleToSync := range roles {
		log.Info().Msg(fmt.Sprintf(
			"sync role %s",
			roleToSync.Name,
		))

		// try and retrieve the role
		retrieveRootRoleResponse, err := roleStore.Retrieve(
			role.RetrieveRequest{
				Filter: filter.NewNameFilter(roleToSync.Name),
			},
		)
		if errors.Is(err, mongo.NewErrNotFound()) {
			// role does not exist yet - create it
			log.Info().Msg("\trole does not exist --> create it")
			if _, err := roleAdmin.AddNewRole(
				role.AddNewRoleRequest{Role: roleToSync},
			); err != nil {
				log.Fatal().Err(err).Msg("\tunable to create role")
			}
			continue
		} else if err != nil {
			// errors other than not "NotFound" are unexpected
			log.Fatal().Err(err).Msg("\terror retrieving role")
		}

		// role does exist, if it has changed update it
		log.Info().Msg("\trole already exists")

		roleToSync.ID = retrieveRootRoleResponse.Role.ID
		if !roleToSync.Equal(retrieveRootRoleResponse.Role) {
			log.Info().Msg("\trole has changed --> update it")
			if _, err := roleAdmin.UpdateRole(
				role.UpdateRoleRequest{Role: roleToSync},
			); err != nil {
				log.Fatal().Err(err).Msg("\tunable to update role")
			}
		}
	}
}
