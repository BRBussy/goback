package role

type Admin interface {
	AddNewRole(AddNewRoleRequest) (*AddNewRoleResponse, error)
	UpdateRole(UpdateRoleRequest) (*UpdateRoleResponse, error)
}

const AdminServiceProviderName = "Role-Admin"

const AddNewRoleService = AdminServiceProviderName + ".AddNewRole"
const UpdateRoleService = AdminServiceProviderName + ".UpdateRole"

type AddNewRoleRequest struct {
	Role Role
}

type AddNewRoleResponse struct {
	Role Role
}

type UpdateRoleRequest struct {
	Role Role
}

type UpdateRoleResponse struct {
	Role Role
}
