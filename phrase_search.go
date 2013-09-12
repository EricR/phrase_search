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
	var records_added int
	var wg sync.WaitGroup
	var words []string
	var word_count int

	// Split phrase into words
	words = strings.Split(entry, " ")
	word_count = len(words)

	// If we don't have any words, don't bother doing anything
	if word_count == 0 {
		return records_added
	}

	// Add our words count to a wait group so we know when we're done
	wg.Add(word_count)

	// Track amunt of time it takes to process the words
	start := time.Now()
	if index.Debug {
		log.Printf("Started writing %d words to '%s'", word_count, index.Name)
	}

	// Create nGrams for all length cases
	if word_count > 1 {
		for bound_size := word_count; bound_size > 0; bound_size-- {
			go func(bound_size int) {
				i_max := word_count - (bound_size - 1)
				for i := 0; i < i_max; i++ {
					phrase := strings.Join(words[i:i+bound_size], " ")
					score := word_count - bound_size
					records_added++
					index.Records[phrase] = append(index.Records[phrase], &Record{data, score})
					if index.Debug {
						log.Printf("%s <-- %s : %s (perm_score=%d)", index.Name, phrase, data, score)
					}
				}

				wg.Done()
			}(bound_size)
		}
	}

	wg.Wait()

	total := time.Now().Sub(start)
	if index.Debug {
		log.Printf("Write took %fs", total.Seconds())
	}

	return records_added
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

	fav_programming_language := index.Find("an awesome programming language")
	fmt.Printf("\nwhat is an awesome programming language? %s (score: %d)\n\n", fav_programming_language[0].Data, fav_programming_language[0].Score)
}
