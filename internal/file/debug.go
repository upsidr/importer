package file

import "fmt"

func (f *File) PrintDebugAll() {
	f.PrintDebugBefore()
	f.PrintDebugPurged()
	f.PrintDebugAll()
}

func (f *File) PrintDebugBefore() {
	fmt.Println("---------------")
	fmt.Println("Content Before:")
	for i, x := range f.contentBefore {
		fmt.Printf("%d: %s\n", i, x)
	}
	fmt.Println("---------------")
	fmt.Println()
}
func (f *File) PrintDebugPurged() {
	fmt.Println("---------------")
	fmt.Println("Content After Purged:")
	for i, x := range f.contentPurged {
		fmt.Printf("%d: %s\n", i, x)
	}
	fmt.Println("---------------")
	fmt.Println()
}
func (f *File) PrintDebugAfter() {
	fmt.Println("---------------")
	fmt.Println("Content After Processed:")
	fmt.Printf("%s", f.contentAfter)
	fmt.Println("---------------")
	fmt.Println()
}
