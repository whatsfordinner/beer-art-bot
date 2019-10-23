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

// SlicesEqual returns true if the two slices have the same elements in the same order, otherwise false.
func SlicesEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// AppendIfUnique is a utility function for only adding a beer or style that has already been seen.
func AppendIfUnique(s []string, e string) []string {
	if len(e) != 0 && !SliceContains(s, e) {
		return append(s, e)
	}
	return s
}

// GetMutuallyExclusiveElements is a utility function that takes in two lists and returns only those
// elements that exist in one list but not both.
func GetMutuallyExclusiveElements(a []string, b []string) []string {
	s := []string{}

	for _, e := range a {
		if !SliceContains(b, e) {
			s = AppendIfUnique(s, e)
		}
	}

	for _, e := range b {
		if !SliceContains(a, e) {
			s = AppendIfUnique(s, e)
		}
	}

	return s
}

// SplitSliceAt splits the input slice into two slices one [0:n] and one [n:len(s)]
func SplitSliceAt(s []string, n int) ([]string, []string) {
	if n >= len(s) {
		return s, []string{}
	}

	if n < 0 {
		return []string{}, s
	}

	return s[:n], s[n:]
}
