package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	verbose bool
	inputFile, filterFiles string
)

func main() {

	files := []string{
		`C:\supp\md5file1`,
		//`C:\supp\md5file2`,
		//`C:\supp\plaintxt1`,
	}

	start := time.Now()
	filters := parseFilters(files);

	fmt.Println("Filters build took ", time.Since(start).Nanoseconds())

	num:=0
	for _, v := range filters.emails {
		num += len(v)
	}
	fmt.Printf("%d emails filter in %d buckets\n", num, len(filters.emails))

	num=0
	for _, v := range filters.md5s {
		num += len(v)
	}
	fmt.Printf("%d dm5s filter in %d buckets\n", num, len(filters.md5s))

	// parse flags
	//ok := parseFlags()
	//if flags !ok
	//print usage

	//ffs :=  strings.Fields(filterFiles)

	// email struct contains email address and optional md5

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
