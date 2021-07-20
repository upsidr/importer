package file

import (
	"fmt"
)

func (f *File) PrintAfter() error {
	fmt.Printf("%s", f.ContentAfter)
	return nil
}

func (f *File) PrintPurged() error {
	fmt.Printf("%s", f.ContentPurged)
	return nil
}
