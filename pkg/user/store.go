package user

import (
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/mongo/filter"
	"os/user"
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
	User user.User `validate:"required"`
}

type CreateResponse struct {
}

type RetrieveRequest struct {
	Filter filter.Filter `validate:"required"`
}

type RetrieveResponse struct {
	User user.User
}

type UpdateRequest struct {
	User user.User `validate:"required"`
}

type UpdateResponse struct {
}

type DeleteRequest struct {
	Filter filter.Filter `validate:"required"`
}

type DeleteResponse struct {
	User user.User
}

type ListRequest struct {
	Filter filter.Filter `validate:"required"`
	Query  mongo.Query
}

type ListResponse struct {
	Records []user.User
	Total   int64
}
