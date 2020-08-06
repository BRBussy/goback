package role

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
	validator  *validate.RequestValidator
	collection *mongo.Collection
}

func NewMongoStore(
	database *mongo.DatabaseConnection,
) Store {
	// get collection
	collection := database.Collection("role")

	// setup collection indices
	if err := collection.SetupIndices(
		[]mongoDriver.IndexModel{
			mongo.NewUniqueIndex("id"),
			mongo.NewUniqueIndex("email"),
		},
	); err != nil {
		log.Fatal().Err(err).Msg("error setting up role collection indices")
	}

	return &MongoStore{
		validator:  validate.NewRequestValidator(),
		collection: collection,
	}
}

func (m MongoStore) Create(request CreateRequest) (*CreateResponse, error) {
	// validate service request
	if err := m.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// insert role entity into mongo collection
	if err := m.collection.CreateOne(request.Role); err != nil {
		log.Error().Err(err).Msg("error creating role")
		return nil, exception.NewErrUnexpected(err)
	}

	return &CreateResponse{}, nil
}

func (m MongoStore) Retrieve(request RetrieveRequest) (*RetrieveResponse, error) {
	// validate service request
	if err := m.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// try and find a role in the collection using the given filter
	var result Role
	if err := m.collection.FindOne(&result, request.Filter); err != nil {
		if errors.Is(err, &mongo.ErrNotFound{}) {
			// if the error is ErrNotFound return it
			return nil, err
		}

		// otherwise return an unexpected error
		log.Error().Err(err).Msg("error retrieving role")
		return nil, exception.NewErrUnexpected(err)
	}

	return &RetrieveResponse{Role: result}, nil
}

func (m MongoStore) Update(request UpdateRequest) (*UpdateResponse, error) {
	// validate service request
	if err := m.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// try and update the role in the collection
	if err := m.collection.UpdateOne(request.Role, filter.NewIDFilter(request.Role.ID)); err != nil {
		// failure here is an unexpected error
		log.Error().Err(err).Msg("error updating role")
		return nil, exception.NewErrUnexpected(err)
	}

	return &UpdateResponse{}, nil
}

func (m MongoStore) Delete(request DeleteRequest) (*DeleteResponse, error) {
	// validate service request
	if err := m.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// try remove a role from the collection using the given filter
	if err := m.collection.DeleteOne(request.Filter); err != nil {
		// failure here is an unexpected error
		log.Error().Err(err).Msg("error deleting role")
		return nil, exception.NewErrUnexpected(err)
	}

	return &DeleteResponse{}, nil
}

func (m MongoStore) List(request ListRequest) (*ListResponse, error) {
	// validate the service request
	if err := m.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// perform find many
	var records []Role
	count, err := m.collection.FindMany(&records, request.Filter, request.Query)
	if err != nil {
		// failure here is an unexpected error
		log.Error().Err(err).Msg("error deleting role")
		return nil, exception.NewErrUnexpected(err)
	}

	return &ListResponse{
		Records: records,
		Total:   count,
	}, nil
}
