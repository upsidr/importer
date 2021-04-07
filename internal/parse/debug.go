package parse

import "fmt"

func (f *File) PrintDebug() {
	fmt.Printf("---------------\n")
	fmt.Printf("Content Before:\n")
	for i, x := range f.contentBefore {
		fmt.Printf("%d: %s\n", i, x)
	}
	fmt.Printf("---------------\n")
	fmt.Println()
	fmt.Printf("---------------\n")
	fmt.Printf("Content After Purged:\n")
	for i, x := range f.contentPurged {
		fmt.Printf("%d: %s\n", i, x)
	}
	fmt.Printf("---------------\n")
	fmt.Println()
	fmt.Printf("---------------\n")
	fmt.Printf("Content After Processed:\n")
	fmt.Printf("%s", f.contentAfter)
	fmt.Printf("---------------\n")
	fmt.Println()
}
