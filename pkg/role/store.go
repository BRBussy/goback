package role

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

const StoreServiceProviderName = "Role-Store"

const CreateService = StoreServiceProviderName + ".Create"
const RetrieveService = StoreServiceProviderName + ".Retrieve"
const UpdateService = StoreServiceProviderName + ".Update"
const DeleteService = StoreServiceProviderName + ".Delete"
const ListService = StoreServiceProviderName + ".List"

type CreateRequest struct {
	Role Role `validate:"required"`
}

type CreateResponse struct {
}

type RetrieveRequest struct {
	Filter filter.Filter `validate:"required"`
}

type RetrieveResponse struct {
	Role Role
}

type UpdateRequest struct {
	Role Role `validate:"required"`
}

type UpdateResponse struct {
}

type DeleteRequest struct {
	Filter filter.Filter `validate:"required"`
}

type DeleteResponse struct {
	Role Role
}

type ListRequest struct {
	Filter filter.Filter `validate:"required"`
	Query  mongo.Query
}

type ListResponse struct {
	Records []Role
	Total   int64
}
