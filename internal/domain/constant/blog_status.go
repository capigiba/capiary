package constant

type BlogStatus string

const (
	BlogStatusActive   BlogStatus = "active"
	BlogStatusInactive BlogStatus = "inactive"
	BlogStatusDeleted  BlogStatus = "deleted"
	BlogStatusArchived BlogStatus = "archived"
)

var AllBlogStatus = []BlogStatus{
	BlogStatusActive,
	BlogStatusInactive,
	BlogStatusDeleted,
	BlogStatusArchived,
}

func IsValidBlogStatus(status BlogStatus) bool {
	return IsValid(status, AllBlogStatus)
}
