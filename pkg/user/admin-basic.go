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
	"golang.org/x/crypto/bcrypt"
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

func (b *BasicAdmin) AddNewUser(request AddNewUserRequest) (*AddNewUserResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// password should be set explicitly
	request.User.Password = make([]byte, 0)

	// registration should be set explicitly
	request.User.Registered = false

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
					Filter: filter.NewTextExactFilter(
						"id",
						roleID,
					),
				},
			); errors.Is(err, &mongo.ErrNotFound{}) {
				reasonsInvalid = append(
					reasonsInvalid,
					fmt.Sprintf("role with ID %s does not exist", roleID),
				)
			} else if err != nil {
				// errors other than not "NotFound" are unexpected
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
				Filter: filter.NewTextExactFilter(
					"email",
					request.User.Email,
				),
			},
		); err == nil {
			// if there was no error during retrieval
			// a user with this email already exists
			reasonsInvalid = append(
				reasonsInvalid,
				"email already in use",
			)
		} else if !errors.Is(err, &mongo.ErrNotFound{}) {
			// errors other than not "NotFound" are unexpected
			log.Error().Err(err).Msg("error retrieving user")
			return nil, exception.NewErrUnexpected(err)
		}
	}

	// confirm username set
	if request.User.Email == "" {
		reasonsInvalid = append(
			reasonsInvalid,
			"username not set",
		)
	} else {
		// if it is set then try and retrieve a user by it
		// to check if the is already in use
		if _, err := b.userStore.Retrieve(
			RetrieveRequest{
				Filter: filter.NewTextExactFilter(
					"username",
					request.User.Username,
				),
			},
		); err == nil {
			// if there was no error during retrieval
			// a user with this username already exists
			reasonsInvalid = append(
				reasonsInvalid,
				"username already in use",
			)
		} else if !errors.Is(err, &mongo.ErrNotFound{}) {
			// errors other than not "NotFound" are unexpected
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

func (b *BasicAdmin) UpdateUser(request UpdateUserRequest) (*UpdateUserResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// try and retrieve the user that is to be updated
	retrieveUserResponse, err := b.userStore.Retrieve(
		RetrieveRequest{
			Filter: filter.NewTextExactFilter(
				"id",
				request.User.ID,
			),
		},
	)
	if errors.Is(err, &mongo.ErrNotFound{}) {
		return nil, NewErrUserDoesNotExist()
	} else if err != nil {
		// errors other than not "NotFound" are unexpected
		log.Error().Err(err).Msg("error retrieving user")
		return nil, exception.NewErrUnexpected(err)
	}

	// password needs to be explicitly set
	request.User.Password = retrieveUserResponse.User.Password

	// confirm that changes were actually made
	if retrieveUserResponse.User.Equal(request.User) {
		return nil, exception.NewErrNoChangesMade()
	}

	// validate the user for update
	reasonsInvalid := make([]string, 0)

	// if the email address has changed confirm that
	// the new email address is not already in use
	if request.User.Email != retrieveUserResponse.User.Email {
		if _, err := b.userStore.Retrieve(
			RetrieveRequest{
				Filter: filter.NewTextExactFilter(
					"email",
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
			// errors other than not "NotFound" are unexpected
			log.Error().Err(err).Msg("error retrieving user")
			return nil, exception.NewErrUnexpected(err)
		}
	}

	// confirm that all of the roles exist
	// then confirm that each is references a valid role
	for _, roleID := range request.User.RoleIDs {
		if _, err := b.roleStore.Retrieve(
			role.RetrieveRequest{
				Filter: filter.NewTextExactFilter(
					"id",
					roleID,
				),
			},
		); errors.Is(err, &mongo.ErrNotFound{}) {
			reasonsInvalid = append(
				reasonsInvalid,
				fmt.Sprintf("role with ID %s is does not exist", roleID),
			)
		} else if err != nil {
			// errors other than not "NotFound" are unexpected
			log.Error().Err(err).Msg("error retrieving role")
			return nil, exception.NewErrUnexpected(err)
		}
	}

	if len(reasonsInvalid) > 0 {
		return nil, NewErrUserNotValid(reasonsInvalid)
	}

	// update the user
	if _, err := b.userStore.Update(
		UpdateRequest{User: request.User},
	); err != nil {
		log.Error().Err(err).Msg("unable to update user")
	}

	return &UpdateUserResponse{}, nil
}

func (b *BasicAdmin) RegisterUser(request RegisterUserRequest) (*RegisterUserResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// try and retrieve the user that is to be registered
	retrieveUserResponse, err := b.userStore.Retrieve(
		RetrieveRequest{
			Filter: filter.NewTextExactFilter(
				"id",
				request.UserID,
			),
		},
	)
	if errors.Is(err, &mongo.ErrNotFound{}) {
		return nil, NewErrUserDoesNotExist()
	} else if err != nil {
		// errors other than not "NotFound" are unexpected
		log.Error().Err(err).Msg("error retrieving user")
		return nil, exception.NewErrUnexpected(err)
	}

	// set the user's password
	if _, err := b.SetUserPassword(
		SetUserPasswordRequest(request),
	); err != nil {
		log.Error().Err(err).Msg("error setting user's password")
		return nil, exception.NewErrUnexpected(err)
	}

	// set the user's registered flag
	retrieveUserResponse.User.Registered = true

	// update the user
	if _, err := b.UpdateUser(
		UpdateUserRequest{User: retrieveUserResponse.User},
	); err != nil {
		log.Error().Err(err).Msg("error updating user")
		return nil, exception.NewErrUnexpected(err)
	}

	return &RegisterUserResponse{}, nil
}

func (b *BasicAdmin) SetUserPassword(request SetUserPasswordRequest) (*SetUserPasswordResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// try and retrieve the user that is to be registered
	retrieveUserResponse, err := b.userStore.Retrieve(
		RetrieveRequest{
			Filter: filter.NewTextExactFilter(
				"id",
				request.UserID,
			),
		},
	)
	if errors.Is(err, &mongo.ErrNotFound{}) {
		return nil, NewErrUserDoesNotExist()
	} else if err != nil {
		// errors other than not "NotFound" are unexpected
		log.Error().Err(err).Msg("error retrieving user")
		return nil, exception.NewErrUnexpected(err)
	}

	// generate a hash of the given password
	pwdHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Error().Err(err).Msg("error hashing password")
		return nil, exception.NewErrUnexpected(err)
	}

	// update the password hash on the user
	retrieveUserResponse.User.Password = pwdHash

	if _, err := b.userStore.Update(
		UpdateRequest{User: retrieveUserResponse.User},
	); err != nil {
		log.Error().Err(err).Msg("error updating user")
		return nil, exception.NewErrUnexpected(err)
	}

	return &SetUserPasswordResponse{}, nil
}

func (b *BasicAdmin) CheckUserPassword(request CheckUserPasswordRequest) (*CheckUserPasswordResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// try and retrieve the user whose password is to be checked
	retrieveUserResponse, err := b.userStore.Retrieve(
		RetrieveRequest{
			Filter: filter.NewTextExactFilter(
				"id",
				request.UserID,
			),
		},
	)
	if errors.Is(err, &mongo.ErrNotFound{}) {
		return nil, NewErrUserDoesNotExist()
	} else if err != nil {
		// errors other than not "NotFound" are unexpected
		log.Error().Err(err).Msg("error retrieving user")
		return nil, exception.NewErrUnexpected(err)
	}

	// check password
	if err := bcrypt.CompareHashAndPassword(
		retrieveUserResponse.User.Password,
		[]byte(request.Password),
	); errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return &CheckUserPasswordResponse{Correct: false}, nil
	} else if err != nil {
		log.Error().Err(err).Msg("error checking password against hash")
		return nil, exception.NewErrUnexpected(err)
	}

	return &CheckUserPasswordResponse{Correct: true}, nil
}
