package user

import (
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/mongo/filter"
	"net/http"
)

type StoreAuthorisedJSONRPCWrapper struct {
	store Store
}

func NewStoreAuthorisedJSONRPCWrapper(
	userStore Store,
) *StoreAuthorisedJSONRPCWrapper {
	return &StoreAuthorisedJSONRPCWrapper{store: userStore}
}

func (a StoreAuthorisedJSONRPCWrapper) ServiceProviderName() string {
	return StoreServiceProviderName
}

type ListJSONRPCRequest struct {
	Filter filter.SerializedFilter `json:"filter"`
	Query  mongo.Query             `json:"query"`
}

type ListJSONRPCResponse struct {
	Records []User `json:"records"`
	Total   int64  `json:"total"`
}

func (a *StoreAuthorisedJSONRPCWrapper) List(r *http.Request, request *ListJSONRPCRequest, response *ListJSONRPCResponse) error {
	listResponse, err := a.store.List(
		ListRequest{
			Filter: request.Filter.Filter,
			Query:  request.Query,
		},
	)
	if err != nil {
		return err
	}

	response.Records = listResponse.Records
	response.Total = listResponse.Total

	return nil
}

type RetrieveJSONRPCRequest struct {
	Filter filter.SerializedFilter `json:"filter"`
}

type RetrieveJSONRPCResponse struct {
	User User `json:"user"`
}

func (a *StoreAuthorisedJSONRPCWrapper) Retrieve(r *http.Request, request *RetrieveJSONRPCRequest, response *RetrieveJSONRPCResponse) error {
	retrieveResponse, err := a.store.Retrieve(
		RetrieveRequest{
			Filter: request.Filter.Filter,
		},
	)
	if err != nil {
		return err
	}

	response.User = retrieveResponse.User

	return nil
}
