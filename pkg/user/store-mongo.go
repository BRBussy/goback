package user

import (
	"errors"
	"github.com/BRBussy/goback/pkg/exception"
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/mongo/filter"
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type MongoStore struct {
	requestValidator *validate.RequestValidator
	collection       *mongo.Collection
}

func NewMongoStore(
	database *mongo.DatabaseConnection,
) Store {
	// get collection
	collection := database.Collection("user")

	// setup collection indices
	if err := collection.SetupIndices(
		[]mongoDriver.IndexModel{
			mongo.NewUniqueIndex("id"),
			mongo.NewUniqueIndex("email"),
		},
	); err != nil {
		log.Fatal().Err(err).Msg("error setting up user collection indices")
	}

	return &MongoStore{
		requestValidator: validate.NewRequestValidator(),
		collection:       collection,
	}
}

func (m MongoStore) Create(request CreateRequest) (*CreateResponse, error) {
	if err := m.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// insert user entity into mongo collection
	if err := m.collection.CreateOne(request.User); err != nil {
		log.Error().Err(err).Msg("error creating user")
		return nil, exception.NewErrUnexpected(err)
	}

	return &CreateResponse{}, nil
}

func (m MongoStore) Retrieve(request RetrieveRequest) (*RetrieveResponse, error) {
	if err := m.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// try and find a user in the collection using the given filter
	var result User
	if err := m.collection.FindOne(&result, request.Filter); err != nil {
		if errors.Is(err, &mongo.ErrNotFound{}) {
			// if the error is ErrNotFound return it
			return nil, err
		}

		// otherwise return an unexpected error
		log.Error().Err(err).Msg("error retrieving user")
		return nil, exception.NewErrUnexpected(err)
	}

	return &RetrieveResponse{User: result}, nil
}

func (m MongoStore) Update(request UpdateRequest) (*UpdateResponse, error) {
	if err := m.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// try and update the user in the collection
	if err := m.collection.UpdateOne(
		request.User,
		filter.NewTextExactFilter(
			"id",
			request.User.ID,
		),
	); err != nil {
		// failure here is an unexpected error
		log.Error().Err(err).Msg("error updating user")
		return nil, exception.NewErrUnexpected(err)
	}

	return &UpdateResponse{}, nil
}

func (m MongoStore) Delete(request DeleteRequest) (*DeleteResponse, error) {
	if err := m.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// try remove a user from the collection using the given filter
	if err := m.collection.DeleteOne(request.Filter); err != nil {
		// failure here is an unexpected error
		log.Error().Err(err).Msg("error deleting user")
		return nil, exception.NewErrUnexpected(err)
	}

	return &DeleteResponse{}, nil
}

func (m MongoStore) List(request ListRequest) (*ListResponse, error) {
	if err := m.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// perform find many
	records := make([]User, 0)
	count, err := m.collection.FindMany(&records, request.Filter, request.Query)
	if err != nil {
		// failure here is an unexpected error
		log.Error().Err(err).Msg("error deleting user")
		return nil, exception.NewErrUnexpected(err)
	}

	return &ListResponse{
		Records: records,
		Total:   count,
	}, nil
}
