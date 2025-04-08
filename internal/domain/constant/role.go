package constant

type Role string

const (
	RoleBasic   Role = "basic"
	RolePremium Role = "premium"
	RoleCP      Role = "cp"
	RoleAdmin   Role = "admin"
)

var AllRoles = []Role{RoleBasic, RolePremium, RoleAdmin, RoleCP}

func IsValidRole(role Role) bool {
	return IsValid(role, AllRoles)
}
