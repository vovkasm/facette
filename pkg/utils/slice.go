package utils

// SliceIndexOfString returns the postion of a string in a slice, or -1 if not found.
func SliceIndexOfString(slice []string, str string) int {
	for i, s := range slice {
		if s == str {
			return i
		}
	}

	return -1
}
