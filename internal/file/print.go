package file

import (
	"fmt"
	"strings"
)

func (f *File) PrintAfter() error {
	fmt.Printf("%s", f.ContentAfter)
	return nil
}

func (f *File) PrintPurged() error {
	s := combineLines(f.ContentPurged)
	fmt.Printf("%s", s)
	return nil
}

func (f *File) PrintBefore() error {
	s := combineLines(f.ContentBefore)
	fmt.Printf("%s", s)
	return nil
}

func combineLines(ss []string) string {
	d := strings.Join(ss, "\n")
	return d
}
