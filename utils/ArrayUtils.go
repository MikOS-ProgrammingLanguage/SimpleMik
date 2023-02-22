package utils

// Returns true if a given string (a) is contained in (b)
func StringInArray(a string, b []string) bool {
	for _, c := range b {
		if c == a {
			return true
		}
	}
	return false
}

// return true is a given int (a) is contained in (b)
func Int8InArray(a int8, b []int8) bool {
	for _, c := range b {
		if c == a {
			return true
		}
	}
	return false
}
