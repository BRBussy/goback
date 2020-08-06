package user

import (
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type MongoStore struct {
	collection       *mongo.Collection
	requestValidator *validate.RequestValidator
}

func NewMongoStore(
	databaseConnection *mongo.DatabaseConnection,
) *MongoStore {
	// get collection
	userCollection := databaseConnection.Collection("user")

	// setup collection indices
	if err := userCollection.SetupIndices(
		[]mongoDriver.IndexModel{
			mongo.NewUniqueIndex("id"),
			mongo.NewUniqueIndex("email"),
		},
	); err != nil {
		log.Fatal().Err(err).Msg("error setting up user collection indices")
	}

	return &MongoStore{
		collection: databaseConnection.Collection("user"),
	}
}

func (s *MongoStore) CreateOne(request CreateOneRequest) (*CreateOneResponse, error) {
	if err := s.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if err := s.collection.CreateOne(request.User); err != nil {
		return nil, NewErrUnexpected(err)
	}
	return &CreateOneResponse{}, nil
}

func (s *MongoStore) FindOne(request FindOneRequest) (*FindOneResponse, error) {
	if err := s.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	var result User
	if err := s.collection.FindOne(&result, request.Identifier); err != nil {
		switch err.(type) {
		case *mongo.ErrNotFound:
			return nil, err
		default:
			return nil, NewErrUnexpected(err)
		}
	}
	return &FindOneResponse{User: result}, nil
}

func (s *MongoStore) FindMany(request FindManyRequest) (*FindManyResponse, error) {
	if err := s.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	var records []User
	count, err := s.collection.FindMany(&records, request.Criteria, request.Query)
	if err != nil {
		return nil, NewErrUnexpected(err)
	}
	if records == nil {
		records = make([]User, 0)
	}

	return &FindManyResponse{
		Records: records,
		Total:   count,
	}, nil
}

func (s *MongoStore) UpdateOne(request UpdateOneRequest) (*UpdateOneResponse, error) {
	if err := s.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if err := s.collection.UpdateOne(request.User, bson.M{}); err != nil {
		return nil, NewErrUnexpected(err)
	}
	return &UpdateOneResponse{}, nil
}
