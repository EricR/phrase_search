package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// A record has associated data and a score.
// Scores are the number of permutations that were required for that record.
// The lower the score, the bigger the nGram token was and the more unique the record is.
//
type Record struct {
	Data  string
	Score int
}

// An index has a name, associated records, and can show debug output.
//
type Index struct {
	Name    string
	Records map[string][]*Record
	Debug   bool
}

// Creates a new index.
//
func NewIndex(name string, debug bool) *Index {
	if debug {
		log.Printf("Creating index '%s'", name)
	}
	return &Index{name, make(map[string][]*Record), debug}
}

// Generates nGram tokens and stores them in the index.
//
func (index *Index) Add(entry string, data string) int {
	var rc int
	var wc int
	var wg sync.WaitGroup
	var words []string

	// Split phrase into words
	words = strings.Split(entry, " ")
	wc = len(words)

	// If word count is 0, we have nothing to do
	if wc == 0 {
		return rc
	}

	// Since we're creating go routines for each iteration of nGram tokenization,
	// we need a waitgroup to know when we're all done
	wg.Add(wc)

	// Track amount of time nGrams tokenization takes
	start := time.Now()
	if index.Debug {
		log.Printf("Started writing %d words to '%s'", wc, index.Name)
	}

	// Generate nGram tokens for all cases of n and create associated records
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

	// Wait until we're finished
	wg.Wait()

	total := time.Now().Sub(start)
	if index.Debug {
		log.Printf("Write took %fs", total.Seconds())
	}

	return wc
}

// Looks up the phrase in an index's records
//
func (index *Index) Find(phrase string) []*Record {
	return index.Records[phrase]
}

func main() {
	index := NewIndex("facts", true)

	fmt.Printf("\nTelling 'facts' that 'my adorable pet dog' = 'Spot'\n\n")
	index.Add("my adorable pet dog", "Spot")

	good_boy := index.Find("pet dog")
	if good_boy != nil {
		fmt.Printf("\nwho is a pet dog? %s (score: %d)\n", good_boy[0].Data, good_boy[0].Score)
	}

	fmt.Printf("\nTelling 'facts' that 'an awesome programming language that i enjoy' = 'Go'\n\n")
	index.Add("an awesome programming language that i enjoy", "Go")

	programming_language := index.Find("an awesome programming language")
	if programming_language != nil {
		fmt.Printf("\nwhat is an awesome programming language? %s (score: %d)\n\n", programming_language[0].Data, programming_language[0].Score)
	}
}
