package file

import (
	"bufio"
	"bytes"
	"fmt"
)

func (f *File) PrintDebugAll() {
	f.PrintDebugBefore()
	f.PrintDebugPurged()
	f.PrintDebugAfter()
}

func (f *File) PrintDebugBefore() {
	fmt.Println("---------------")
	fmt.Println("Content Before:")
	for i, x := range f.ContentBefore {
		fmt.Printf("%d:\t%s\n", i, x)
	}
	fmt.Println("---------------")
	fmt.Println()
}
func (f *File) PrintDebugPurged() {
	fmt.Println("---------------")
	fmt.Println("Content After Purged:")
	for i, x := range f.ContentPurged {
		fmt.Printf("%d:\t%s\n", i, x)
	}
	fmt.Println("---------------")
	fmt.Println()
}
func (f *File) PrintDebugAfter() {
	fmt.Println("---------------")
	fmt.Println("Content After Processed:")
	currentLine := 0
	scanner := bufio.NewScanner(bytes.NewReader(f.ContentAfter))
	for scanner.Scan() {
		currentLine++
		fmt.Printf("%d:\t%s\n", currentLine, scanner.Text())
	}
	fmt.Println("---------------")
	fmt.Println()
}
