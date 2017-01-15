package main

import (
	"flag"
	"fmt"
	"strings"
	"./classifying"
)

var (
	verbose bool
	inputFile, filterFiles string
)

func main() {
	if ok := parseFlags(); ok {
		ffs := strings.Fields(filterFiles)
		filter := supress.ParseFilters(ffs);
		supress.ProcessInput(inputFile, filter);
	}
	fmt.Println("Wrong usage. Try --help")
}

func parseFlags() bool {
	flag.BoolVar(&verbose, "verbose", false, "prompt more information about the working progress")
	flag.StringVar(&inputFile, "input", "", "input emails filename")
	flag.StringVar(&filterFiles, "filters", "", "filter filenames list. Should be double-quoted whitespace separated")
	flag.Parse()

	return inputFile != "" && filterFiles != ""
}
