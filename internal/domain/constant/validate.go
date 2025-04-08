package constant

func IsValid[T comparable](val T, list []T) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}
