package user

type Admin interface {
	AddNewUser(AddNewUserRequest) (*AddNewUserResponse, error)
}

const AdminServiceProviderName = "User-Admin"

const AddNewUserService = AdminServiceProviderName + ".AddNewUser"

type AddNewUserRequest struct {
	User User
}

type AddNewUserResponse struct {
	User User
}
