package user

import (
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/security/claims"
)

type Store interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	FindOne(FindOneRequest) (*FindOneResponse, error)
	FindMany(FindManyRequest) (*FindManyResponse, error)
	UpdateOne(UpdateOneRequest) (*UpdateOneResponse, error)
}

const ServiceProvider = "User-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const FindOneService = ServiceProvider + ".FindOne"
const FindManyService = ServiceProvider + ".FindMany"
const UpdateOneService = ServiceProvider + ".UpdateOne"

type CreateOneRequest struct {
	User User `validate:"required"`
}

type CreateOneResponse struct {
}

type FindOneRequest struct {
	Claims     claims.Claims `validate:"required"`
	Identifier mongo.Filter  `validate:"required"`
}

type FindOneResponse struct {
	User User
}

type FindManyRequest struct {
	Claims   claims.Claims `validate:"required"`
	Criteria mongo.Filter  `validate:"required"`
	Query    mongo.Query
}

type FindManyResponse struct {
	Records []User
	Total   int64
}

type UpdateOneRequest struct {
	Claims claims.Claims `validate:"required"`
	User   User
}

type UpdateOneResponse struct {
}
