package user

import (
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/mongo/filter"
)

type Store interface {
	Create(CreateRequest) (*CreateResponse, error)
	Retrieve(RetrieveRequest) (*RetrieveResponse, error)
	Update(UpdateRequest) (*UpdateResponse, error)
	Delete(DeleteRequest) (*DeleteResponse, error)
	List(ListRequest) (*ListResponse, error)
}

const ServiceProvider = "User-Store"

const CreateService = ServiceProvider + ".Create"
const RetrieveService = ServiceProvider + ".Retrieve"
const UpdateService = ServiceProvider + ".Update"
const DeleteService = ServiceProvider + ".Delete"
const ListService = ServiceProvider + ".List"

type CreateRequest struct {
	User User `validate:"required"`
}

type CreateResponse struct {
}

type RetrieveRequest struct {
	Filter filter.Filter `validate:"required"`
}

type RetrieveResponse struct {
	User User
}

type UpdateRequest struct {
	User User `validate:"required"`
}

type UpdateResponse struct {
}

type DeleteRequest struct {
	Filter filter.Filter `validate:"required"`
}

type DeleteResponse struct {
	User User
}

type ListRequest struct {
	Filter filter.Filter `validate:"required"`
	Query  mongo.Query
}

type ListResponse struct {
	Records []User
	Total   int64
}
