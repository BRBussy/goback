package user

import (
	"errors"
	"fmt"
	"github.com/BRBussy/goback/pkg/exception"
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/mongo/filter"
	"github.com/BRBussy/goback/pkg/role"
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
	"strings"
)

type BasicAdmin struct {
	requestValidator *validate.RequestValidator
	userStore        Store
	roleStore        role.Store
}

func NewBasicAdmin(
	userStore Store,
	roleStore role.Store,
) *BasicAdmin {
	return &BasicAdmin{
		requestValidator: validate.NewRequestValidator(),
		userStore:        userStore,
		roleStore:        roleStore,
	}
}

func (b BasicAdmin) AddNewUser(request AddNewUserRequest) (*AddNewUserResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// clean up user entity
	request.User.Password = make([]byte, 0)
	request.User.Email = strings.ReplaceAll(strings.ToLower(request.User.Email), " ", "")

	// validate the user for creation
	reasonsInvalid := make([]string, 0)

	// if the user has any assigned role IDs
	if request.User.RoleIDs == nil {
		request.User.RoleIDs = make([]string, 0)
	} else {
		// then confirm that each is references a valid role
		for _, roleID := range request.User.RoleIDs {
			if _, err := b.roleStore.Retrieve(
				role.RetrieveRequest{
					Filter: filter.NewIDFilter(roleID),
				},
			); errors.Is(err, &mongo.ErrNotFound{}) {
				reasonsInvalid = append(
					reasonsInvalid,
					fmt.Sprintf("role with ID %s does not exist", roleID),
				)
			} else if err != nil {
				// if there was an error that is not "NotFound" this is an unexpected error
				log.Error().Err(err).Msg("error retrieving role")
				return nil, exception.NewErrUnexpected(err)
			}
		}
	}

	// confirm email set
	if request.User.Email == "" {
		reasonsInvalid = append(
			reasonsInvalid,
			"email not set",
		)
	} else {
		// if it is set then try and retrieve a user by it
		// to check if the is already in use
		if _, err := b.userStore.Retrieve(
			RetrieveRequest{
				Filter: filter.NewEmailFilter(
					request.User.Email,
				),
			},
		); err == nil {
			// if there was no error during retrieval
			// a user with this email address already exists
			reasonsInvalid = append(
				reasonsInvalid,
				"email already in use",
			)
		} else if !errors.Is(err, &mongo.ErrNotFound{}) {
			// if there was an error that is not "NotFound" this is an unexpected error
			log.Error().Err(err).Msg("error retrieving user")
			return nil, exception.NewErrUnexpected(err)
		}
	}

	if len(reasonsInvalid) > 0 {
		return nil, NewErrUserNotValid(reasonsInvalid)
	}

	// set user ID
	request.User.ID = uuid.NewV4().String()

	// create user
	if _, err := b.userStore.Create(
		CreateRequest{User: request.User},
	); err != nil {
		log.Error().Err(err).Msg("unable to create user")
		return nil, exception.NewErrUnexpected(err)
	}

	return &AddNewUserResponse{User: request.User}, nil
}

func (b BasicAdmin) UpdateUser(request UpdateUserRequest) (*UpdateUserResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &UpdateUserResponse{}, nil
}
