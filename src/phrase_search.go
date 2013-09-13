package main

import (
	"fmt"
	"github.com/streadway/simpleuuid"
	"log"
	"strings"
	"sync"
	"time"
)

type Token struct {
	Document *Document
	Score    float32
}

type Document struct {
	UUID simpleuuid.UUID
	Data string
}

type Index struct {
	Name   string
	Tokens map[string][]*Token
	Debug  bool
}

type SearchResult struct {
	TokenRef *Token
	Document *Document
}

func (sr *SearchResult) Score() float32 {
	return sr.TokenRef.Score
}

func NewDocument(data string) *Document {
	uuid, _ := simpleuuid.NewTime(time.Now())
	return &Document{uuid, data}
}

func NewIndex(name string, debug bool) *Index {
	if debug {
		log.Printf("Creating index '%s'", name)
	}
	return &Index{name, make(map[string][]*Token), debug}
}

func (index *Index) Insert(text string, data string) int {
	var wg sync.WaitGroup
	var phrases []string

	wcounter := 0
	document := NewDocument(data)
	token_map := make(map[string]*Token)
	phrases = strings.FieldsFunc(text, SentenceDelims)
	time_start := time.Now()

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

					token_map[phrase] = &Token{document, score}
					if index.Debug {
						log.Printf("%s <-- %s : %s (perm_score=%d)", index.Name, phrase, data, score)
					}
				}

				wg.Done()
			}(n)
		}

		wg.Wait()
	}

	time_total := time.Now().Sub(time_start)
	log.Printf("Wrote %d words to '%s' (took %fs)", wcounter, index.Name, time_total.Seconds())

	for p, r := range token_map {
		index.Tokens[p] = append(index.Tokens[p], r)
	}

	return len(token_map)
}

func (index *Index) Search(phrase string) []*Token {
	return index.Tokens[strings.ToLower(phrase)]
}

func SentenceDelims(r rune) bool {
	return r == '.' || r == ',' || r == '?' || r == '!'
}

func main() {
	index := NewIndex("Lorem Ipsums", false)

	fmt.Printf("\nWriting lorem ipsums to 'facts'\n\n")
	index.Insert("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut metus varius laoreet. Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue. Curabitur ullamcorper ultricies nisi. A needle in a hay stack", "First Lorem Ipsum")
	index.Insert("Sed lectus. Integer euismod lacus luctus magna. Quisque cursus, metus vitae pharetra auctor, sem massa mattis sem, at interdum magna augue eget diam. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Morbi lacinia molestie dui. Praesent blandit dolor. Sed non quam. In vel mi sit amet augue congue elementum. Morbi in ipsum sit amet pede facilisis laoreet. Donec lacus nunc, viverra nec, blandit vel, egestas et, augue. Vestibulum tincidunt malesuada tellus. Ut ultrices ultrices enim.", "Second Loreum Ipsum")
	index.Insert("Mauris ipsum. Nulla metus metus, ullamcorper vel, tincidunt sed, euismod in, nibh. Quisque volutpat condimentum velit. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Nam nec ante. Sed lacinia, urna non tincidunt mattis, tortor neque adipiscing diam, a cursus ipsum ante quis turpis. Nulla facilisi. Ut fringilla. Suspendisse potenti. Nunc feugiat mi a tellus consequat imperdiet. Vestibulum sapien. Proin quam. Etiam ultrices. Suspendisse in justo eu magna luctus suscipit.", "Third Loreum Ipsum")
	index.Insert("Curabitur sit amet mauris. Morbi in dui quis est pulvinar ullamcorper. Nulla facilisi. Integer lacinia sollicitudin massa. Cras metus. Sed aliquet risus a tortor. Integer id quam. Morbi mi. Quisque nisl felis, venenatis tristique, dignissim in, ultrices sit amet, augue. Proin sodales libero eget ante. Nulla quam. Aenean laoreet. Vestibulum nisi lectus, commodo ac, facilisis ac, ultricies eu, pede. Ut orci risus, accumsan porttitor, cursus quis, aliquet eget, justo. ", "Fourth Loreum Ispum")

	fmt.Printf("\nSearching for 'a needle in a Hay Stack' in 'Lorem Ipsums' index")
	time_start := time.Now()
	needle := index.Search("a needle in a Hay Stack")
	time_total := time.Now().Sub(time_start).Seconds()

	if needle != nil {
		fmt.Printf("\nFound %d document(s) in %fs: Document{uuid: %s, data: \"%s\", score: %1.2f}\n\n", len(needle), time_total, needle[0].Document.UUID, needle[0].Document.Data, needle[0].Score)
	}
}
