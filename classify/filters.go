package supress

import (
	"regexp"
	"sync"
)

const (
	splittersNum = 5
	keyLength = 3
)

var md5Regex = regexp.MustCompile(`^[a-f0-9]{32}$`)

type multiMap map[string][]string

type filters struct {
	emails     multiMap
	md5s       multiMap
	md5Enabled bool
}

// reads filter files provided splits to md5 and emails and
// returns filter struct
func ParseFilters(files []string) *filters {
	var f = &filters{
		emails     : make(multiMap),
		md5s       : make(multiMap),
		md5Enabled : false,
	}

	for _, n := range files {
		parseFilter(n, f)
	}
	f.md5Enabled = len(f.md5s) != 0

	return f
}

func parseFilter(file string, f*filters) {
	//run reader that will produce chunks to next pipeline consumer
	lines := stream(file)

	emails := make(chan string, 10000)
	md5s := make(chan string, 10000)
	//run classifiers
	var wg sync.WaitGroup
	wg.Add(splittersNum)
	for w := 0; w < splittersNum; w++ {
		go split(lines, emails, md5s, &wg)
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

func split(lines <-chan string, emails, md5s chan <- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for l := range lines {
		if md5Regex.MatchString(l) {
			md5s <- l
		} else {
			emails <- l
		}
	}
}

func fill(f*filters, emails, md5s <-chan string) {
	var wg sync.WaitGroup

	fillMap := func(mm*multiMap, data <-chan string) {
		defer wg.Done()
		for d := range data {
			k := d[0:keyLength]
			(*mm)[k] = append((*mm)[k], d)
		}
	}

	wg.Add(2)
	go fillMap(&f.emails, emails)
	go fillMap(&f.md5s, md5s)

	wg.Wait()
}
