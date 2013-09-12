package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type Record struct {
	Data  string
	Score int
}

type Index struct {
	Name    string
	Records map[string][]*Record
	Debug   bool
}

func NewIndex(name string, debug bool) *Index {
	if debug {
		log.Printf("Creating index '%s'", name)
	}
	return &Index{name, make(map[string][]*Record), debug}
}

func (index *Index) Add(entry string, data string) int {
	var rc int
	var wc int
	var wg sync.WaitGroup
	var words []string

	// Split phrase into words
	words = strings.Split(entry, " ")
	wc = len(words)

	// If word count is zero, we have nothing to do
	if wc == 0 {
		return rc
	}

	// Add our words count to a wait group so we have all nGrams
	wg.Add(wc)

	// Track amount of time nGram generation takes
	start := time.Now()
	if index.Debug {
		log.Printf("Started writing %d words to '%s'", wc, index.Name)
	}

	// Create nGrams for all cases of n
	for n := wc; n > 0; n-- {
		go func(n int) {
			i_max := wc - (n - 1)
			for i := 0; i < i_max; i++ {
				phrase := strings.Join(words[i:i+n], " ")
				score := wc - n
				rc++

				index.Records[phrase] = append(index.Records[phrase], &Record{data, score})
				if index.Debug {
					log.Printf("%s <-- %s : %s (perm_score=%d)", index.Name, phrase, data, score)
				}
			}

			wg.Done()
		}(n)
	}

	// Wait until we have all nGrams
	wg.Wait()

	total := time.Now().Sub(start)
	if index.Debug {
		log.Printf("Write took %fs", total.Seconds())
	}

	return wc
}

func (index *Index) Find(phrase string) []*Record {
	return index.Records[phrase]
}

func main() {
	index := NewIndex("facts", true)

	fmt.Printf("\nTelling 'facts' that 'my adorable pet dog' = 'Spot'\n\n")
	index.Add("my adorable pet dog", "Spot")

	good_boy := index.Find("pet dog")
	fmt.Printf("\nwho is a pet dog? %s (score: %d)\n", good_boy[0].Data, good_boy[0].Score)

	fmt.Printf("\nTelling 'facts' that 'an awesome programming language that i enjoy' = 'Go'\n\n")
	index.Add("an awesome programming language that i enjoy", "Go")

	programming_language := index.Find("an awesome programming language")
	fmt.Printf("\nwhat is an awesome programming language? %s (score: %d)\n\n", programming_language[0].Data, programming_language[0].Score)
}
