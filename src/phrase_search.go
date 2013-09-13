package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// A record has associated data and a score.
// Scores are a percentage of phrase permutations over total sentence permutations
//
type Record struct {
	Data  string
	Score float32
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

	// Split text into phrases
	phrases = strings.FieldsFunc(text, SentenceDelims)

	start := time.Now()

	for _, phrase := range phrases {
    phrase = strings.ToLower(phrase)
		words := strings.Split(phrase, " ")
		wc := len(words)

		if wc == 0 {
			break
		}

		wg.Add(wc)

		for n := wc; n > 0; n-- {
			go func(n int) {
				i_max := wc - (n - 1)
				for i := 0; i < i_max; i++ {
					phrase := strings.Join(words[i:i+n], " ")
					score := float32(n) / float32(wc)
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
	log.Printf("Wrote %d words to '%s' (took %fs)", wcounter, index.Name, total.Seconds())

	for p, r := range records {
		index.Records[p] = append(index.Records[p], r)
	}

	return len(records)
}

// Looks up the phrase in an index's records
//
func (index *Index) Search(phrase string) []*Record {
	return index.Records[strings.ToLower(phrase)]
}

// Sentnece delimiters
//
func SentenceDelims(r rune) bool {
	return r == '.' || r == ',' || r == '?' || r == '!'
}

func main() {
	index := NewIndex("Lorem Ipsums", false)

	fmt.Printf("\nWriting lorem ipsums to 'facts'\n\n")
	index.Add("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut metus varius laoreet. Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue. Curabitur ullamcorper ultricies nisi. A needle in a hay stack", "First Lorem Ipsum")
	index.Add("Sed lectus. Integer euismod lacus luctus magna. Quisque cursus, metus vitae pharetra auctor, sem massa mattis sem, at interdum magna augue eget diam. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Morbi lacinia molestie dui. Praesent blandit dolor. Sed non quam. In vel mi sit amet augue congue elementum. Morbi in ipsum sit amet pede facilisis laoreet. Donec lacus nunc, viverra nec, blandit vel, egestas et, augue. Vestibulum tincidunt malesuada tellus. Ut ultrices ultrices enim.", "Second Loreum Ipsum")
	index.Add("Mauris ipsum. Nulla metus metus, ullamcorper vel, tincidunt sed, euismod in, nibh. Quisque volutpat condimentum velit. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Nam nec ante. Sed lacinia, urna non tincidunt mattis, tortor neque adipiscing diam, a cursus ipsum ante quis turpis. Nulla facilisi. Ut fringilla. Suspendisse potenti. Nunc feugiat mi a tellus consequat imperdiet. Vestibulum sapien. Proin quam. Etiam ultrices. Suspendisse in justo eu magna luctus suscipit.", "Third Loreum Ipsum")
	index.Add("Curabitur sit amet mauris. Morbi in dui quis est pulvinar ullamcorper. Nulla facilisi. Integer lacinia sollicitudin massa. Cras metus. Sed aliquet risus a tortor. Integer id quam. Morbi mi. Quisque nisl felis, venenatis tristique, dignissim in, ultrices sit amet, augue. Proin sodales libero eget ante. Nulla quam. Aenean laoreet. Vestibulum nisi lectus, commodo ac, facilisis ac, ultricies eu, pede. Ut orci risus, accumsan porttitor, cursus quis, aliquet eget, justo. ", "Fourth Loreum Ispum")

  fmt.Printf("\nSearching for 'a needle in a Hay Stack' in 'Lorem Ipsums' index")
  time_start := time.Now()
  needle := index.Search("a needle in a Hay Stack")
  time_total := time.Now().Sub(time_start).Seconds()

  if needle != nil {
    fmt.Printf("\nFound %d record(s) in %fs: %s (score: %1.2f)\n\n", len(needle), time_total, needle[0].Data, needle[0].Score)
  }
}
