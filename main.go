package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	verbose bool
	inputFile, filterFiles string
)

func main() {

	// parse flags
	ok := parseFlags()
	//if flags !ok
		//print usage

	filters :=  strings.Fields(filterFiles)


	// start log

	fmt.Println(filters)

	// emails[] = read input file into multimap with 2 first letters as key[]. The key is basic split point for go-routines
	// filters[first_two_symbols, [emails]], md5_filters[] = read all filters and merge/split

	// for each email in emails
		//run filter:
			//if not is_valid_mail
				//add to bad
			//else if exists_in_filters
				//add to matches
			//else if calculate md5 exists_in_filters
				//add to matches
			//else
				//add to clean

	// dump collected results using buffers.
}

func parseFlags() bool {
	flag.BoolVar(&verbose, "verbose", false, "prompt more information about the working progress")
	flag.StringVar(&inputFile, "input", "", "input emails filename")
	flag.StringVar(&filterFiles, "filters", "", "filter filenames")
	flag.Parse()

	return inputFile != "" && filterFiles != ""
}