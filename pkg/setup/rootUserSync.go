package setup

import (
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/mongo/filter"
	"github.com/BRBussy/goback/pkg/role"
	"github.com/BRBussy/goback/pkg/user"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var rootUser = user.User{
	Username: "root",
	Email:    "root@example.com",
	RoleIDs:  make([]string, 0),
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
		},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("error listing root user roles")
	}
	if len(listRolesResponse.Records) != len(rootUserRoleNames) {
		log.Fatal().Msg("incorrect number of roles listed")
	}

	// populate all of the roleIDs on the root user
	for _, r := range listRolesResponse.Records {
		rootUser.RoleIDs = append(
			rootUser.RoleIDs,
			r.ID,
		)
	}

	// try and retrieve the root user by username
	retrieveRootUserResponse, err := userStore.Retrieve(
		user.RetrieveRequest{
			Filter: filter.NewUsernameFilter(rootUser.Username),
		},
	)
	switch err.(type) {
	case *mongo.ErrNotFound:
		// user does not exist yet - create and register it

		log.Info().Msg("root user does not exist")
		log.Info().Msg("\t--> create it")
		addNewUserResponse, err := userAdmin.AddNewUser(
			user.AddNewUserRequest{User: rootUser},
		)
		if err != nil {
			log.Fatal().Err(err).Msg("error adding root user")
		}
		log.Info().Msg("\t--> register it")
		if _, err := userAdmin.RegisterUser(
			user.RegisterUserRequest{
				UserID:   addNewUserResponse.User.ID,
				Password: rootPassword,
			},
		); err != nil {
			log.Fatal().Err(err).Msg("error registering user")
		}

	case nil:
		// user already exists - update if required
		log.Info().Msg("root user already exists")

		// generate a hash of the root password
		pwdHash, err := bcrypt.GenerateFromPassword(
			[]byte(rootPassword),
			bcrypt.DefaultCost,
		)
		if err != nil {
			log.Fatal().Err(err).Msg("error hashing password")
		}

		// set pwd and id on root user entity
		rootUser.Password = pwdHash
		rootUser.ID = retrieveRootUserResponse.User.ID
		rootUser.Registered = true

		// check if update required
		if rootUser.Equal(retrieveRootUserResponse.User) {
			log.Info().Msg("\t--> no changes")
			return
		}
		log.Info().Msg("\t--> changes made")
		if _, err := userAdmin.UpdateUser(
			user.UpdateUserRequest{User: rootUser},
		); err != nil {
			log.Fatal().Err(err).Msg("error updating root user")
		}

	default:
		// errors other than not "NotFound" are unexpected
		log.Fatal().Err(err).Msg("error retrieving root user")
	}
}
