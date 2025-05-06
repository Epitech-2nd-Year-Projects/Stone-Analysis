package main

import (
	"flag"
	"fmt"
	"os"
	"stone-analysis/internal/analyze"
	"stone-analysis/internal/cypher"
	"stone-analysis/internal/decypher"
	"stone-analysis/internal/utils"
	"strconv"
)

func main() {
	analyzeFlag := flag.Bool("analyze", false, "Run in analyze mode")
	cypherFlag := flag.Bool("cypher", false, "Run in cypher mode")
	decypherFlag := flag.Bool("decypher", false, "Run in decypher mode")

	flag.Parse()

	args := flag.Args()

	modesSet := 0
	if *analyzeFlag {
		modesSet++
	}
	if *cypherFlag {
		modesSet++
	}
	if *decypherFlag {
		modesSet++
	}

	if modesSet == 0 {
		utils.DisplayHelp()
		os.Exit(84)
	}

	if modesSet > 1 {
		utils.DisplayHelp()
		os.Exit(84)
	}

	if *analyzeFlag {
		if len(args) != 2 {
			utils.DisplayHelp()
			os.Exit(84)
		}

		inFile := args[0]
		nStr := args[1]

		if err := utils.CheckFileExists(inFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.DisplayHelp()
			os.Exit(84)
		}

		n, err := strconv.Atoi(nStr)
		if err != nil {
			utils.DisplayHelp()
			os.Exit(84)
		}

		if n < 1 {
			utils.DisplayHelp()
			os.Exit(84)
		}
		analyze.Analyze(inFile, n)
	} else if *cypherFlag {
		if len(args) != 3 {
			utils.DisplayHelp()
			os.Exit(84)
		}

		inFile := args[0]
		outFile := args[1]
		message := args[2]

		if err := utils.CheckFileExists(inFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.DisplayHelp()
			os.Exit(84)
		}
		cypher.Cypher(inFile, outFile, message)
	} else if *decypherFlag {
		if len(args) != 1 {
			utils.DisplayHelp()
			os.Exit(84)
		}

		inFile := args[0]

		if err := utils.CheckFileExists(inFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.DisplayHelp()
			os.Exit(84)
		}
		decypher.Decypher(inFile)
	}
}
