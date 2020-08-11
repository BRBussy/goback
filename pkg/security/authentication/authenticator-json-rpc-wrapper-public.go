package authentication

import (
	"net/http"
)

type AuthenticatorAuthorisedJSONRPCWrapper struct {
	authenticator Authenticator
}

func NewAuthenticatorAuthorisedJSONRPCWrapper(authenticator Authenticator) *AuthenticatorAuthorisedJSONRPCWrapper {
	return &AuthenticatorAuthorisedJSONRPCWrapper{authenticator: authenticator}
}

func (a AuthenticatorAuthorisedJSONRPCWrapper) ServiceProviderName() string {
	return AuthenticatorServiceProviderName
}

type LoginJSONRPCRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginJSONRPCResponse struct {
	JWT string `json:"jwt"`
}

func (a *AuthenticatorAuthorisedJSONRPCWrapper) Login(r *http.Request, request *LoginJSONRPCRequest, response *LoginJSONRPCResponse) error {
	addNewUserResponse, err := a.authenticator.Login(
		LoginRequest(*request),
	)
	if err != nil {
		return err
	}

	response.JWT = addNewUserResponse.JWT

	return nil
}
