package sliceutil

// SliceContains is a utility function for determining if a beer or style has been seen already.
func SliceContains(s []string, e string) bool {
	for _, i := range s {
		if i == e {
			return true
		}
	}
	return false
}

// AppendIfUnique is a utility function for only adding a beer or style that has already been seen.
func AppendIfUnique(s []string, e string) []string {
	if len(e) != 0 && !SliceContains(s, e) {
		return append(s, e)
	}
	return s
}
