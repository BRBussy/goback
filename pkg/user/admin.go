package user

type Admin interface {
	Get(GetRequest) (*GetResponse, error)
	Set(SetRequest) (*SetResponse, error)
}

const AdminServiceProviderName = "User-Admin"

const GetService = AdminServiceProviderName + ".Get"
const SetService = AdminServiceProviderName + ".Set"

type GetRequest struct {
	User User
}

type GetResponse struct {
	User User
}

type SetRequest struct {
	User User
}

type SetResponse struct {
	User User
}
