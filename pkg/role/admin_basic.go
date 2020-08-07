package role

import (
	"errors"
	"github.com/BRBussy/goback/pkg/exception"
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/mongo/filter"
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
)

type BasicAdmin struct {
	requestValidator *validate.RequestValidator
	roleStore        Store
}

func NewBasicAdmin(
	roleStore Store,
) *BasicAdmin {
	return &BasicAdmin{
		requestValidator: validate.NewRequestValidator(),
		roleStore:        roleStore,
	}
}

func (b BasicAdmin) AddNewRole(request AddNewRoleRequest) (*AddNewRoleResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// validate the new role for creation
	if request.Role.Name == "" {
		return nil, NewErrRoleNotValid([]string{"name not set"})
	}

	// try and retrieve the role by name to check if it already exists
	if _, err := b.roleStore.Retrieve(
		RetrieveRequest{
			Filter: filter.NewNameFilter(
				request.Role.Name,
			),
		},
	); err == nil {
		// if there was no error during retrieval, the role already exists
		return nil, NewErrRoleAlreadyExists()
	} else if !errors.Is(err, &mongo.ErrNotFound{}) {
		// if there was an error that is not "NotFound" this is an unexpected error
		log.Error().Err(err).Msg("error retrieving role")
		return nil, exception.NewErrUnexpected(err)
	}

	// role is valid and can be created

	// set ID and permissions fields
	request.Role.ID = uuid.NewV4().String()
	if request.Role.Permissions == nil {
		request.Role.Permissions = make([]string, 0)
	}

	if _, err := b.roleStore.Create(
		CreateRequest{
			Role: request.Role,
		},
	); err != nil {
		log.Error().Err(err).Msg("error creating role")
		return nil, exception.NewErrUnexpected(err)
	}

	return &AddNewRoleResponse{}, nil
}

func (b BasicAdmin) UpdateRole(request UpdateRoleRequest) (*UpdateRoleResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &UpdateRoleResponse{}, nil
}
