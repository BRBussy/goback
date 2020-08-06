package role

type Admin interface {
	Get(GetRequest) (*GetResponse, error)
	Set(SetRequest) (*SetResponse, error)
}

const AdminServiceProviderName = "Role-Admin"

const GetService = AdminServiceProviderName + ".Get"
const SetService = AdminServiceProviderName + ".Set"

type GetRequest struct {
	Role Role
}

type GetResponse struct {
	Role Role
}

type SetRequest struct {
	Role Role
}

type SetResponse struct {
	Role Role
}
