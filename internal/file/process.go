package file

import (
	"fmt"
)

const br = byte('\n')

// ProcessMarkers reads markers and generates ContentAfter.
//
// Internally, ContentAfter is generated from ContentPurged and Markers.
// This walks thruogh each line of ContentPurged, and copies the data into a
// byte slice, and while processing each line, it checks if markers are
// defined for the given line. If any marker is registered, it would then
// process the target information to import.
//
// TODO: possibly remove error return, as it currently never returns any error.
func (f *File) ProcessMarkers() error {
	result := []byte{}
	for line, data := range f.ContentPurged {
		result = append(result, data...)
		result = append(result, br)

		// Marker is found for the given line. Before proceeding to the
		// next line, handle marker and import the target data.
		if marker, found := f.Markers[line+1]; found {
			processed, err := marker.ProcessMarkerData(f.FileName)
			if err != nil {
				fmt.Printf("warning: %s\n", err)
				continue
			}
			result = append(result, processed...)
		}
	}
	f.ContentAfter = result
	return nil
}
