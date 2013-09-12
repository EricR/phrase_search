package main

import (
	"fmt"
	"log"
	"strings"
	"time"
  "sync"
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

func (index *Index) Add(phrase string, data string) int {
	words := strings.Split(phrase, " ")
	wc := len(words)
	var records_added int
	var nw string
	var ns int
  var wg sync.WaitGroup

	if wc == 0 {
		return records_added
	}

  wg.Add(wc)

	start := time.Now()
	if index.Debug {
		log.Printf("Started writing %d words to '%s'", wc, index.Name)
	}

	if wc > 1 {
		for bs := wc; bs > 0; bs-- {
			go func(bs int) {
				for i := 0; i < wc-(bs-1); i++ {
					nw = strings.Join(words[i:i+bs], " ")
					ns = wc - bs
					records_added++
					index.Records[nw] = append(index.Records[nw], &Record{data, ns})
					if index.Debug {
						log.Printf("%s <-- %s : %s (perm_score=%d)", index.Name, nw, data, ns)
					}
				}

        wg.Done()
			}(bs)
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

	fav_programming_language := index.Find("awesome programming language")
	fmt.Printf("\nwhat is an awesome programming language? %s (score: %d)\n\n", fav_programming_language[0].Data, fav_programming_language[0].Score)
}
