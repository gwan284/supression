package main

import (
	"flag"
	"fmt"
	"time"
	"./classifying"
)

var (
	verbose bool
	inputFile, filterFiles string
)

func main() {

	files := []string{
		//`C:\supp\md5file1`,
		//`C:\supp\md5file2`,
		`C:\supp\plaintxt1`,
	}

	start := time.Now()
	filters := supress.ParseFilters(files);

	fmt.Println("Filters build took ", time.Since(start).Nanoseconds(), &filters)


	// parse flags
	//ok := parseFlags()
	//if flags !ok
	//print usage

	//ffs :=  strings.Fields(filterFiles)

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
