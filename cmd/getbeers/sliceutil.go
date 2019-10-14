package main

// sliceContains is a utility function for determining if a beer or style has been seen already.
func sliceContains(s []string, e string) bool {
	for _, i := range s {
		if i == e {
			return true
		}
	}
	return false
}

// appendIfUnique is a utility function for only adding a beer or style that has already been seen.
func appendIfUnique(s []string, e string) []string {
	if len(e) != 0 && !sliceContains(s, e) {
		return append(s, e)
	}
	return s
}
