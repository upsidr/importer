package parse

import (
	"bufio"
	"strings"
	"testing"
)

// StringToLineBytes converts input string into slice of byte slices, where
// each outer slice element represents a single line.
func StringToLineBytes(t testing.TB, data string) [][]byte {
	t.Helper()

	result := make([][]byte, 0)
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		result = append(result, scanner.Bytes())
	}

	return result
}
