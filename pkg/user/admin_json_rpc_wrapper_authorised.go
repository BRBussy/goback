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

type AddNewUserJSONRPCRequest struct {
	User User `json:"user"`
}

type AddNewUserJSONRPCResponse struct {
	User User `json:"user"`
}

func (a *AdminAuthorisedJSONRPCWrapper) AddNewUser(r *http.Request, request *AddNewUserJSONRPCRequest, response *AddNewUserJSONRPCResponse) error {
	addNewUserResponse, err := a.admin.AddNewUser(
		AddNewUserRequest(*request),
	)
	if err != nil {
		return err
	}

	response.User = addNewUserResponse.User

	return err
}
