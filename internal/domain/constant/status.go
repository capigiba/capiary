package constant

type AccountStatus string

const (
	StatusActive    AccountStatus = "active"
	StatusInactive  AccountStatus = "inactive"
	StatusPending   AccountStatus = "pending"
	StatusSuspended AccountStatus = "suspended"
	StatusBanned    AccountStatus = "banned"
	StatusDeleted   AccountStatus = "deleted"
	StatusArchived  AccountStatus = "archived"
)

var AllAccountStatus = []AccountStatus{
	StatusActive,
	StatusInactive,
	StatusPending,
	StatusSuspended,
	StatusBanned,
	StatusDeleted,
	StatusArchived,
}

func IsValidAccountStatus(accountStatus AccountStatus) bool {
	for _, r := range AllAccountStatus {
		if r == accountStatus {
			return true
		}
	}
	return false
}
