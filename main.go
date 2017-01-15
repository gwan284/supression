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

	input := `C:\supp\emailfile`

	start := time.Now()
	fs := supress.ParseFilters(files);
	end := time.Now()
	fmt.Println("Filters build took ", time.Since(start).Seconds())

	supress.ProcessInput(input, fs);

	fmt.Println("Sorting input according filter took ", time.Since(end).Seconds())


	// parse flags
	//ok := parseFlags()
	//if flags !ok
	//print usage

	//ffs :=  strings.Fields(filterFiles)
}

func parseFlags() bool {
	flag.BoolVar(&verbose, "verbose", false, "prompt more information about the working progress")
	flag.StringVar(&inputFile, "input", "", "input emails filename")
	flag.StringVar(&filterFiles, "filters", "", "filter filenames")
	flag.Parse()

	return inputFile != "" && filterFiles != ""
}
