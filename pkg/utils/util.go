package utils

// If a given string in a slice.
func InSlice(haystack []string, niddle string) bool {
	for _, v := range haystack {
		if v == niddle {
			return true
		}
	}

	return false
}
