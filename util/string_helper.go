package util

func IsBlank(val string) bool {
	if len(val) == 0 {
		return true
	}
	return false
}

func IsAnyBlank(val ...string) bool {
	for _, s := range val {
		if IsBlank(s) {
			return true
		}
	}
	return false
}
