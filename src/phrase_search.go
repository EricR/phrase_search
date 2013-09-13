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
func (index *Index) Add(text string, data string) int {
  var wg sync.WaitGroup
  var phrases []string
  var wcounter int

  // Init records
  records := make(map[string]*Record)

  // Make text all lowercase
  text = strings.ToLower(text)

	// Split text into phrases
  phrases = strings.Split(text, ".")
  phrases = strings.Split(text, "!")
  phrases = strings.Split(text, "?")
  phrases = strings.Split(text, ",")

  start := time.Now()
  log.Printf("Started writing %d phrases to '%s'", len(phrases), index.Name)

  for _, phrase := range phrases {
    words := strings.Split(phrase, " ")
	  wc := len(words)

    // If word count is 0, we have nothing to do
    if wc == 0 {
      break
    }

    wg.Add(wc)

    // Generate nGram tokens for all cases of n and create associated records
    for n := wc; n > 0; n-- {
      go func(n int) {
        i_max := wc - (n - 1)
        for i := 0; i < i_max; i++ {
          phrase := strings.Join(words[i:i+n], " ")
          score := wc - n
          wcounter++

          records[phrase] = &Record{data, score}
          if index.Debug {
            log.Printf("%s <-- %s : %s (perm_score=%d)", index.Name, phrase, data, score)
          }
        }

        wg.Done()
      }(n)
    }

    wg.Wait()
  }

  total := time.Now().Sub(start)
  log.Printf("Wrote %d words (%fs)", wcounter, total.Seconds())

  for p, r := range records {
    index.Records[p] = append(index.Records[p], r)
  }

	return len(records)
}

// Looks up the phrase in an index's records
//
func (index *Index) Find(phrase string) []*Record {
	return index.Records[strings.ToLower(phrase)]
}

func main() {
	index := NewIndex("facts", false)

  fmt.Printf("\nWriting lorem ipsum to 'facts'\n\n")
  index.Add("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut metus varius laoreet. Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue. Curabitur ullamcorper ultricies nisi. A needle in a hay stack", "It works!")
  
  needle := index.Find("a needle in a Hay Stack")
  if needle != nil {
    fmt.Printf("\nFinding 'a needle in a Hay Stack' Found %d record(s). Returned data: %s (score: %d)\n\n", len(needle), needle[0].Data, needle[0].Score)
  }
}
