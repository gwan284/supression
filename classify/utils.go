package supress

import (
	"os"
	"log"
	"bufio"
	"sync"
)

func stream(file string) <-chan string {
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

func write(file string, lines <-chan string, wg *sync.WaitGroup, counter *int) {
	defer wg.Done()

	f, err := os.OpenFile(file, os.O_CREATE | os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for l := range lines {
		_, err := f.WriteString(l + "\n")
		(*counter)++
		if err != nil {
			log.Fatal(err)
		}
	}
}

func contains(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}
