package marker

import (
	"bytes"
	"strings"
)

func adjustIndentation(lineData []byte, exporterMarkerIndent int, importerIndentation *Indentation) []byte {
	// If no indentation setup is done, simply return as is
	if importerIndentation == nil {
		lineData = append(lineData, br)
		return lineData
	}

	// Check which indentation adjustment is used.
	// Absolute adjustment takes precedence over extra indentation.
	switch importerIndentation.Mode {
	case AbsoluteIndentation:
		lineData = handleAbsoluteIndentation(lineData, exporterMarkerIndent, importerIndentation.Length)
	case ExtraIndentation:
		lineData = prependWhitespaces(lineData, importerIndentation.Length)
	case AlignIndentation:
		lineData = handleAbsoluteIndentation(lineData, exporterMarkerIndent, importerIndentation.MarkerIndentation)
	case KeepIndentation: // Explicitly handling this, as it is likely that the default behaviour woulld need to change
	}
	lineData = append(lineData, br)
	return lineData
}

// handleAbsoluteIndentation updates the lineData with provided indentation.
//
// There are 3 different indent information to handle:
//
//   A. Indent target in absolute value
//   B. Original preceding indent in the lineData
//   C. Indent of Exporter Marker
//
// With the absolute indentation handling, A and B are obviously required. For
// C, it is important that YAML tree structure is persisted, and thus it is
// crucial to know how much indent Exporter Marker has.
//
// With the above covered, there are a few cases to handle:
//
//   1. Original lineData has more indentation than target indentation
//   2. Original lineData has fewer indentation than target indentation
//   3. Original lineData has fewer preceding indentation than Importer Marker
//
// For the Case 1. and 2., the diff needs to be calculated to ensure correct
// indentation.
//
// For the Case 3., it's not clear what we should expect. This is currently not
// handled, and it may need to be an error.
func handleAbsoluteIndentation(lineData []byte, exportMarkerIndent, targetIndent int) []byte {
	lineString := string(lineData)
	currenttIndent := len(lineString) - len(strings.TrimLeft(lineString, " "))

	switch {
	// Case 1.
	// Marker appears with more indentation than Absolute, and thus strip
	// extra indentations.
	case exportMarkerIndent >= targetIndent:
		indentAdjustment := exportMarkerIndent - targetIndent
		return lineData[indentAdjustment:]

	// Case 2.
	// Marker has less indentation than Absolute wants, and thus prepend
	// the indent diff.
	case exportMarkerIndent < targetIndent:
		indentAdjustment := targetIndent - exportMarkerIndent
		return prependWhitespaces(lineData, indentAdjustment)

	// Case 3.
	case currenttIndent < exportMarkerIndent:
		// TODO: Handle case where indentation is less than marker indentation
	}
	return lineData
}

func prependWhitespaces(x []byte, count int) []byte {
	empty := bytes.Repeat([]byte(" "), count)
	// x = append(x, empty...)
	// copy(x[count:], x)
	// copy(x, empty)
	// return x

	return append(empty, x...)
}
