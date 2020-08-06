package user

import (
	"net/http"
)

type AdminAuthorisedJSONRPCWrapper struct {
	admin Admin
}

func NewAdminAuthorisedJSONRPCWrapper(admin Admin) *AdminAuthorisedJSONRPCWrapper {
	return &AdminAuthorisedJSONRPCWrapper{admin: admin}
}

func (a AdminAuthorisedJSONRPCWrapper) ServiceProviderName() string {
	return AdminServiceProviderName
}

type GetJSONRPCRequest struct {
	User User `json:"user"`
}

type GetJSONRPCResponse struct {
	User User `json:"user"`
}

func (a *AdminAuthorisedJSONRPCWrapper) Get(r *http.Request, request *GetJSONRPCRequest, response *GetJSONRPCResponse) error {
	getResponse, err := a.admin.Get(
		GetRequest(*request),
	)
	if err != nil {
		return err
	}
	response.User = getResponse.User
	return err
}
