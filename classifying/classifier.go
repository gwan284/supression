package supress

import (
	"regexp"
	"sync"
	"crypto/md5"
	"encoding/hex"
)

var emailRegex = regexp.MustCompile(`^(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})$`)

const classifiersNum = 10

func ProcessInput(file string, f *filters) {
	lines := stream(file)

	matches, clean, bad := classify(lines, f)

	//collect data into single combined filter using blocking read channels
	save(matches, clean, bad, file)
}

func classify(ls <-chan string, f *filters) (<-chan string, <-chan string, <-chan string) {
	matches := make(chan string, 10000)
	clean := make(chan string, 10000)
	bad := make(chan string, 10000)

	var wg sync.WaitGroup
	wg.Add(classifiersNum)
	for w := 0; w < classifiersNum; w++ {
		go runClassifier(ls, f, matches, clean, bad, &wg)
	}

	go func() {
		wg.Wait()
		close(matches)
		close(clean)
		close(bad)
	}()

	return matches, clean, bad
}

func runClassifier(ls <-chan string, f *filters, m, c, b chan <- string, wg *sync.WaitGroup) {
	defer wg.Done()

	matchesFilter := func(l string, m multiMap, kLength int) bool {
		if v, ok := m[l[0:kLength]]; ok {
			return contains(v, l)
		}
		return false
	}

	matchesMd5 := func(l string, m multiMap) bool {
		md5 := md5.Sum([]byte(l))
		md5str := hex.EncodeToString(md5[:])

		return matchesFilter(md5str, m, md5KeyLength)
	}

	for l := range ls {
		//if not is_valid_mail send to bad
		if !emailRegex.MatchString(l) {
			b <- l
			//if exists_in_filters send to matches
		} else if matchesFilter(l, f.emails, emailKeyLength) {
			m <- l
			//if calculate md5 exists_in_filters send to matches
		} else if f.md5Enabled && matchesMd5(l, f.md5s) {
			m <- l
			// if non of above send to clean
		} else {
			c <- l
		}
	}
}

func save(m, c, b <-chan string, f string) {
	var wg sync.WaitGroup

	matches, clean, bad := f + ".matches", f + ".clean", f + ".bad"

	wg.Add(3)
	go write(matches, m, &wg)
	go write(clean, c, &wg)
	go write(bad, b, &wg)

	wg.Wait()
}
