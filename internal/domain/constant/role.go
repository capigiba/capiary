package constant

type Role string

const (
	RoleBasic   Role = "basic"
	RolePremium Role = "premium"
	RoleCJ      Role = "cj"
	RoleAdmin   Role = "admin"
)

var AllRoles = []Role{
	RoleBasic,
	RolePremium,
	RoleAdmin,
	RoleCJ,
}

func IsValidRole(role Role) bool {
	return IsValid(role, AllRoles)
}
