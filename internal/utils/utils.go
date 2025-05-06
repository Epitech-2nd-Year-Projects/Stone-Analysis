package utils

import (
	"fmt"
	"os"
)

func DisplayHelp() {
	fmt.Fprintf(
		os.Stdout,
		"USAGE\n%s [--analyze IN_FILE N | --cypher IN_FILE OUT_FILE MESSAGE | --decypher IN_FILE]\n\n",
		os.Args[0],
	)
	fmt.Println("\tIN_FILE\tAn audio file to be analyzed")
	fmt.Println("\tOUT_FILE\tOutput audio file of the cypher mode")
	fmt.Println("\tMESSAGE\tThe message to hide in the audio file")
	fmt.Println("\tN\tNumber of top frequencies to display")
}

func CheckFileExists(filePath string) error {
	fileInfo, err := os.Stat(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("error: input file '%s' does not exist", filePath)
		}
		return fmt.Errorf(
			"error: could not stat input file '%s': %v",
			filePath,
			err,
		)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf(
			"error: input file '%s' is a directory, not a file",
			filePath,
		)
	}
	return nil
}
