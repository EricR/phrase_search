package main

import (
	"fmt"
	"github.com/streadway/simpleuuid"
	"log"
	"strings"
	"sync"
	"time"
)

type TokenCollection []*Token
type TokenMap map[string]*Token
type TokenCollectionMap map[string][]*Token
type DocumentCollectionMap map[*simpleuuid.UUID]*Document

type Token struct {
	Phrase   string
	Document *Document
	Score    float32
}

type Document struct {
	UUID   simpleuuid.UUID
	Index  *Index
	Data   string
	Tokens TokenCollection
}

func (document *Document) Delete() {
	document.Tokens = NewTokenCollection()
	delete(document.Index.Documents, &document.UUID)
}

type Index struct {
	Name      string
	Debug     bool
	Tokens    TokenCollectionMap
	Documents DocumentCollectionMap
}

func (index *Index) Insert(text string, data string) int {
	var wg sync.WaitGroup
	var phrases []string

	document := NewDocument(index, data)
	token_map := make(TokenMap)
	phrases = strings.FieldsFunc(text, SentenceDelims)
	wcounter := 0
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

					token_map[phrase] = &Token{phrase, document, score}
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

	for p, t := range token_map {
		index.Tokens[p] = append(index.Tokens[p], t)
		document.Tokens = append(document.Tokens, t)
	}

	index.Documents[&document.UUID] = document

	return len(token_map)
}

func (index *Index) Search(phrase string) []*Token {
	return index.Tokens[strings.ToLower(phrase)]
}

type SearchResult struct {
	Token
	*Document
}

func (search_result *SearchResult) Score() float32 {
	return search_result.Token.Score
}

func NewTokenCollection() TokenCollection {
	return make(TokenCollection, 100)
}

func NewDocument(index *Index, data string) *Document {
	uuid, _ := simpleuuid.NewTime(time.Now())
	return &Document{uuid, index, data, NewTokenCollection()}
}

func NewIndex(name string, debug bool) *Index {
	if debug {
		log.Printf("Creating index '%s'", name)
	}
	return &Index{name, debug, make(TokenCollectionMap), make(DocumentCollectionMap)}
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
	results := index.Search("a needle in a Hay Stack")
	needle := results[0]
	time_total := time.Now().Sub(time_start).Seconds()
	fmt.Printf("\nFound %d document(s) in %fs: Document{uuid: %s, data: \"%s\", score: %1.2f}\n", len(results), time_total, needle.Document.UUID, needle.Document.Data, needle.Score)

	fmt.Printf("\nIndex now has %d documents and %d tokens", len(index.Documents), len(index.Tokens))

	fmt.Printf("\n\nDeleting found record\n\n")
	needle.Document.Delete()
	fmt.Printf("Index now has %d documents and %d tokens\n\n", len(index.Documents), len(index.Tokens))
}
