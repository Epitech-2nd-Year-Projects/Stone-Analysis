package analyze

import "fmt"

func Analyze(inFile string, n int) error {
	fmt.Printf("Mode: Analyze\n")
	fmt.Printf("IN_FILE: %s\n", inFile)
	fmt.Printf("N: %d\n", n)
	return nil
}
