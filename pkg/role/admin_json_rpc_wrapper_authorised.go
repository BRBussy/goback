package role

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
	Role Role `json:"role"`
}

type GetJSONRPCResponse struct {
	Role Role `json:"role"`
}

func (a *AdminAuthorisedJSONRPCWrapper) Get(r *http.Request, request *GetJSONRPCRequest, response *GetJSONRPCResponse) error {
	getResponse, err := a.admin.Get(
		GetRequest(*request),
	)
	if err != nil {
		return err
	}
	response.Role = getResponse.Role
	return err
}
