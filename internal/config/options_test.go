package config

import "testing"

func TestIgnoreProject(t *testing.T) {

	// create a list of tests to carry out
	testCases := []struct {
		ignore   string
		patterns []string
		expected []bool
	}{
		{
			"testing",
			[]string{"testing"},
			[]bool{true},
		},
		{
			"webapi,crqs",
			[]string{"webapi", "events"},
			[]bool{true, false},
		},
		{
			"webapi",
			[]string{"web.*?", ".*api"},
			[]bool{true, true},
		},
	}

	// iterate around the test cases and perform each test
	for _, testCase := range testCases {

		// create a new options object to work with
		option := Options{}
		option.Ignore = testCase.ignore

		for idx, pattern := range testCase.patterns {
			result := option.IgnoreProject(pattern)

			if result != testCase.expected[idx] {
				t.Errorf("Result of '%s' in '%s' does not equal %v", pattern, testCase.ignore, testCase.expected[idx])
			}
		}

	}
}
