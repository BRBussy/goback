package setup

import (
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/mongo/filter"
	"github.com/BRBussy/goback/pkg/role"
	"github.com/BRBussy/goback/pkg/user"
	"github.com/rs/zerolog/log"
)

var rootUser = user.User{
	Email: "root@goback.com",
}

var rootUserRoleNames = []string{
	"Root",
}

func RootUserSync(
	roleStore role.Store,
	userStore user.Store,
	userAdmin user.Admin,
	rootPassword string,
) {
	// search for all of the root user roles
	listRolesResponse, err := roleStore.List(
		role.ListRequest{
			Filter: filter.NewExactTextListFilter(
				"name",
				rootUserRoleNames,
			),
			Query: mongo.Query{},
		},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("error listing root user roles")
	}
	if len(listRolesResponse.Records) != len(rootUserRoleNames) {
		log.Fatal().Msg("incorrect number of roles listed")
	}
}
