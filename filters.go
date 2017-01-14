package main

import (
	"regexp"
	"bufio"
	"os"
	"log"
	"sync"
)

var md5Regex = regexp.MustCompile(`^[a-f0-9]{32}$`)

type multiMap map[string][]string

type filters struct {
	emails multiMap
	md5s   multiMap
}

func parseFilters(files []string) *filters {
	var f = &filters {
		emails : make(multiMap),
		md5s   : make(multiMap),
	}

	for _, n := range files {
		parseFilter(n, f)
	}

	return f
}

func parseFilter(file string, f* filters) {
	//run reader that will produce chunks to next pipeline consumer
	lines := read(file)

	emails := make(chan string, 10000)
	md5s := make(chan string, 10000)
	//run 5 classifiers
	const classifiersNum = 5
	var wg sync.WaitGroup
	wg.Add(classifiersNum)
	for w := 0; w < classifiersNum; w++ {
		go classify(lines, emails, md5s, &wg)
	}
	go func() {
		wg.Wait()
		close(emails)
		close(md5s)
	}()

	//add optional merge() step to process few files in the same time
	//merge(channels) into one

	//collect data into single combined filter using blocking read channels
	fill(f, emails, md5s)
}

func read(file string) <-chan string {
	lines := make(chan string, 10000)
	go func() {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	return lines
}

func classify(lines <-chan string, emails, md5s chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for l := range lines {
		if md5Regex.MatchString(l) {
			md5s <- l
		} else {
			emails <- l
		}
	}
}

func fill(f* filters, emails, md5s <-chan string) {
	var wg sync.WaitGroup

	fillMap := func(mm* multiMap, data <-chan string) {
		defer wg.Done()
		for d := range data {
			k := d[0:2]
			entry := (*mm)[k]
			(*mm)[k] = append(entry, d)
		}
	}

	wg.Add(2)
	go fillMap(&f.emails, emails)
	go fillMap(&f.md5s, md5s)

	wg.Wait()
}
