package main

import (
	"flag"
	"fmt"
	"log"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"./classify"
)

const filenameLog  = "/var/log/suppression.log"

var (
	optionVerbose bool
	optionInput, optionFilters string
)

var usageFunc =  func() {
		fmt.Fprintf(os.Stderr, "Usage: ./suppress [-v] -input=<input file> -filters=\"<filter file1> [filter file2] [filter file3] ...\"\n")
		flag.PrintDefaults()
	}

func main() {
	flag.Usage = usageFunc
	parseFlags();

	if optionInput == "" || optionFilters == "" {
		fmt.Fprintf(os.Stderr, "No input/filters provided.\n")
		flag.Usage()
		return
	}

	logger := log.New(ioutil.Discard, "", log.LstdFlags)
	if optionVerbose {
		setOutput(logger)
	}

	if _, err := os.Stat(optionInput); os.IsNotExist(err) {
		logger.Fatalf("OpenFile: %v (%v)", err, optionInput)
	}

	var filterFiles []string
	for _, f := range strings.Fields(optionFilters) {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			logger.Printf("Filter %v skipped. OpenFile: %v\n", f, err)
		} else {
			filterFiles = append(filterFiles, f)
		}
	}
	if len(filterFiles) == 0 {
		logger.Fatalln("No valid filters provided. Finishing...")
	}

	logger.Print("Suppression started. Creating filters...")

	filter := supress.ParseFilters(filterFiles);
	logger.Print("Done. Starting input classification...")
	feedback := supress.ProcessInput(optionInput, filter);

	for fb := range feedback {
		logger.Printf("%s matches: %d non-matches: %d\n", optionInput, fb.Matches, fb.Bad + fb.Clean)
	}

	logger.Println("Done!")
}

func parseFlags() {
	flag.BoolVar(&optionVerbose, "v", false, "prompt more information about the working progress")
	flag.StringVar(&optionInput, "input", "", "input emails filename")
	flag.StringVar(&optionFilters, "filters", "", "filter filenames list. Should be double-quoted whitespace separated")
	flag.Parse()
}

func setOutput(l *log.Logger) {
	logfile, err := os.Create(filenameLog)
	if err != nil {
		defer l.Printf("Failed to create log file: %s", err.Error())
	}

	mw := io.MultiWriter(os.Stdout, logfile)
	l.SetOutput(mw)
}
