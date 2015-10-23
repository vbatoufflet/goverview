package slice

// Int64IndexOf returns the index of a given integer in a slice, or -1 if not found
func Int64IndexOf(slice []int64, value int64) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}

	return -1
}

// StringIndexOf returns the index of a given string in a slice, or -1 if not found
func StringIndexOf(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}

	return -1
}
