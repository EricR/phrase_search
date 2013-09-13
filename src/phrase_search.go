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
	log.Printf("Started writing %d words to '%s'", wc, index.Name)

	// Generate nGram tokens for all cases of n and create associated records
	for n := wc; n > 0; n-- {
		go func(n int) {
			i_max := wc - (n - 1)
      records := make(map[string]*Record)
      for i := 0; i < i_max; i++ {
				phrase := strings.Join(words[i:i+n], " ")
				score := wc - n
				rc++

				records[phrase] = &Record{data, score}
				if index.Debug {
					log.Printf("%s <-- %s : %s (perm_score=%d)", index.Name, phrase, data, score)
				}
			}

      for p, r := range records {
        index.Records[p] = append(index.Records[p], r)
      }

			wg.Done()
		}(n)
	}

	// Wait until we're finished
	wg.Wait()

	total := time.Now().Sub(start)
	log.Printf("Write took %fs", total.Seconds())

	return wc
}

// Looks up the phrase in an index's records
//
func (index *Index) Find(phrase string) []*Record {
	return index.Records[phrase]
}

func main() {
	index := NewIndex("facts", false)

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

  fmt.Printf("\nTelling 'facts' a bunch of random stuff\n\n")
  index.Add("this is a test with repeating words repeating words repeating words hopefully it performs well and it doesnt break lets hope so test test test test repeating words repeating words test repeating words", "works")

  test := index.Find("repeating words")
  if test != nil {
    fmt.Printf("\nrepeating words ? %s (score: %d)\n\n", test[0].Data, test[0].Score)
  }

  fmt.Printf("\nTelling 'facts' a lot of stuff\n\n")
  index.Add("ropy up tel pale khayal periphrases a hi mishap periosteorrhaphy hypostoma oh weepy satrapies proatheist uh trilletto a of usure applauses statutorily hulk at saltest asp was yarry ketyl urali rate a hays slitwork to psykters oleo pleomastia lap aortoptosia ye auxiliarly reiter woo oily yuk palmilla a riser haphtarah uprush of laws ply housemaster maltworm ahoy hosel tom thruway my mesolite tomosis lama rye profitless papaprelatist eke uh ketway throwwort furmety mia a multiflue serratiform rosser ha keratometry fopship postholes fly we sup olepy afforest olio host kisra seels oh prutah yip masterwort allorrhyhmia pall rillow hi polythely weaselwise sax pot fatal soporiferous uh up a oafs uppop misappropriates purity why of sap flex elfwife asset so err tits littermates hurt rams rule peal pyrophile tams them me ye upstares pow homoiousious oomph myropolist a toe pulleys ritely frothy khalifas ow petal toe islot tosser uh teras spy phi empresa a extremum this loftless a misstop port a smokeshaft hysteropathy yolk photomappe miss smithite you phyla limitless wholly lustres rex plea hetairas a a hopperette sparse assaut frass swum phloem twaes retypes um part retromammary ye proller oestriasis fart up sootlike impresari pip amyxorrhoea isotypes faitery a maksoorah paw rosy arty malaperts puss emissaries prexy solutes lithemia flatfoots pitau a us trap florae aft lasty surrealists so superearthly ow samel matripotestal slippier harp laius aim alulet skimos septole slaty to tea rokee a far realities sows pre of firmware yep prowfishes uptwist why frithwork imperf me upstep so aah stipitiform arm omit hark sirky ext awes hysteroptosia spermolysis spirket awoke am ha pram emplastrum shat a needle in a hay stack hale sea sat a", "it works")
  
  needle := index.Find("a needle in a hay stack")
  if needle != nil {
    fmt.Printf("\nfound needle in haystack? %s (score: %d)\n\n", needle[0].Data, needle[0].Score)
  }
}
