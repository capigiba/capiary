package constant

type Role string

const (
	RoleBasic   Role = "basic"
	RolePremium Role = "premium"
	RoleAdmin   Role = "admin"
)

var AllRoles = []Role{RoleBasic, RolePremium, RoleAdmin}

func IsValidRole(role Role) bool {
	for _, r := range AllRoles {
		if r == role {
			return true
		}
	}
	return false
}
