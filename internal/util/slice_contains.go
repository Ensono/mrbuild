package util

import "regexp"

func SliceContains(slice []string, value string) bool {
	var result bool

	for _, x := range slice {

		// use the value as a regular expression pattern to use
		// to match against the items in the slice
		re := regexp.MustCompile(value)

		if re.Match([]byte(x)) {
			result = true
			break
		}
	}

	return result
}
