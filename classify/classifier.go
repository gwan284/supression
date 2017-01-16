package supress

import (
	"regexp"
	"sync"
	"crypto/md5"
	"encoding/hex"
	"time"
)

const classifiersNum = 10

var (
	feedbackPeriod = time.Minute
	emailRegex = regexp.MustCompile(`^(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})$`)
)

type feedback struct {
	Matches, Clean, Bad int
}

// runs input emails classification based on filter provided
// returns feedback channel with progress info inside
func ProcessInput(file string, f *filters) (<-chan feedback) {
	lines := stream(file)

	matches, clean, bad := classify(lines, f)

	//collect data into single combined filter using blocking read channels
	 return save(matches, clean, bad, file)
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

// classifier worker routine
// takes emails stream and splits into matches/clean/bad according to rules
func runClassifier(ls <-chan string, f *filters, m, c, b chan <- string, wg *sync.WaitGroup) {
	defer wg.Done()

	matchesFilter := func(l string, m multiMap) bool {
		if v, ok := m[l[0:keyLength]]; ok {
			return contains(v, l)
		}
		return false
	}

	matchesMd5 := func(l string, m multiMap) bool {
		md5 := md5.Sum([]byte(l))
		md5str := hex.EncodeToString(md5[:])

		return matchesFilter(md5str, m)
	}

	for l := range ls {
		//if not is_valid_mail send to bad
		if !emailRegex.MatchString(l) {
			b <- l
			//if exists_in_filters send to matches
		} else if matchesFilter(l, f.emails) {
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

func save(m, c, b <-chan string, f string) (<-chan feedback) {
	fb := make(chan feedback)
	var wg sync.WaitGroup

	matches, clean, bad := f + ".matches", f + ".clean", f + ".bad"
	matchesCnt, cleanCnt, badCnt := 0, 0, 0

	wg.Add(3)
	go write(matches, m, &wg, &matchesCnt)
	go write(clean, c, &wg, &cleanCnt)
	go write(bad, b, &wg, &badCnt)

	// report feedback for logging each feedbackPeriod
	ticker := time.NewTicker(feedbackPeriod)
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				fb <- feedback{ Matches : matchesCnt, Clean : cleanCnt, Bad: badCnt}
			case <- stop:
				ticker.Stop()
				fb <- feedback{ Matches : matchesCnt, Clean : cleanCnt, Bad: badCnt}
				close(fb)
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(stop)
	}()

	return fb
}
