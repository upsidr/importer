package parse

import (
	"bufio"
	"strings"
	"testing"
)

// StringToLineStrings converts input string into slice of byte slices, where
// each outer slice element represents a single line.
func StringToLineStrings(t testing.TB, data string) []string {
	t.Helper()

	result := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result
}
