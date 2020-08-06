package user

import (
	"net/http"
)

type AuthenticatedJSONRPCWrapper struct {
	store Store
}

func NewAuthenticatedJSONRPCWrapper(userStore Store) *AuthenticatedJSONRPCWrapper {
	return &AuthenticatedJSONRPCWrapper{store: userStore}
}

func (a AuthenticatedJSONRPCWrapper) ServiceProviderName() string {
	return StoreServiceProviderName
}

type CreateJSONRPCRequest struct {
	User User `json:"user"`
}

type CreateJSONRPCResponse struct {
}

func (a *AuthenticatedJSONRPCWrapper) Create(r *http.Request, request *CreateJSONRPCRequest, response *CreateJSONRPCResponse) error {
	_, err := a.store.Create(
		CreateRequest(*request),
	)
	return err
}
