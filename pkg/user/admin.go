package user

type Admin interface {
	AddNewUser(AddNewUserRequest) (*AddNewUserResponse, error)
	UpdateUser(UpdateUserRequest) (*UpdateUserResponse, error)
	RegisterUser(RegisterUserRequest) (*RegisterUserResponse, error)
	SetUserPassword(SetUserPasswordRequest) (*SetUserPasswordResponse, error)
}

const AdminServiceProviderName = "User-Admin"

const AddNewUserService = AdminServiceProviderName + ".AddNewUser"
const UpdateUserService = AdminServiceProviderName + ".UpdateUser"

type AddNewUserRequest struct {
	User User `validate:"-"`
}

type AddNewUserResponse struct {
	User User
}

type UpdateUserRequest struct {
	User User `validate:"required"`
}

type UpdateUserResponse struct {
}

type RegisterUserRequest struct {
	UserID   string `validate:"required"`
	Password string `validate:"required"`
}

type RegisterUserResponse struct {
}

type SetUserPasswordRequest struct {
	UserID   string `validate:"required"`
	Password string `validate:"required"`
}

type SetUserPasswordResponse struct {
}
