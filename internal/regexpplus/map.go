package regexpplus

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrNoMatch = errors.New("no match")
)

// MapWithNamedSubgroups runs regexp FindStringSubmatch against `targetLine`
// input using `expression`, and returns a map representation. The map contains
// the key as the subgroup name, and value for the matched data.
//
// This can be used for regular expression which does not have any subgroup,
// but as it is designed specifically for subgroup based use cases, it will
// create a map that will not have all the matched components.
func MapWithNamedSubgroups(targetLine string, expression string) (map[string]string, error) {
	re := regexp.MustCompile(expression)
	ms := re.FindStringSubmatch(targetLine)

	if len(ms) == 0 {
		return nil, fmt.Errorf("%w", ErrNoMatch)
	}

	result := map[string]string{}
	for idx, submatchNames := range re.SubexpNames() {
		result[submatchNames] = ms[idx]
	}

	return result, nil
}
