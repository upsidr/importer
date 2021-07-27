package file

import (
	"bufio"
	"bytes"
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

func (f *File) PrintDebugAll() {
	f.PrintDebugBefore()
	f.PrintDebugPurged()
	f.PrintDebugAfter()
}

func (f *File) PrintDebugBefore() {
	defer wrapWithDivider("Content Before:")()

	for i, x := range f.ContentBefore {
		fmt.Printf("%d:\t%s\n", i, x)
	}
}
func (f *File) PrintDebugPurged() {
	defer wrapWithDivider("Content After Purged:")()

	for i, x := range f.ContentPurged {
		fmt.Printf("%d:\t%s\n", i, x)
	}
}
func (f *File) PrintDebugAfter() {
	defer wrapWithDivider("Content After Processed:")()

	currentLine := 0
	scanner := bufio.NewScanner(bytes.NewReader(f.ContentAfter))
	for scanner.Scan() {
		currentLine++
		fmt.Printf("%d:\t%s\n", currentLine, scanner.Text())
	}
}

func wrapWithDivider(title string) func() {
	divider := "---------------------------------------"

	fmt.Println(divider)
	fmt.Println(title)

	return func() {
		fmt.Println(divider)
		fmt.Println()
	}
}
