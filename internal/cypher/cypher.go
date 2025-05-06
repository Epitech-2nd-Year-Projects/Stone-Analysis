package cypher

import "fmt"

func Cypher(inFile, outFile, message string) error {
	fmt.Printf("Mode: Cypher\n")
	fmt.Printf("IN_FILE: %s\n", inFile)
	fmt.Printf("OUT_FILE: %s\n", outFile)
	fmt.Printf("MESSAGE: %s\n", message)
	return nil
}
