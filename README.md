# Phrase Search

A small text searcher that is optimized to work with bodies of text.
It:

* Takes a paragraph and breaks it up into sentences using punctuation.
* Tokenizes each sentence with a lowercase and nGrams tokenizer.
* Stores the tokens in a records map, which makes sure a single body of
  text won't create duplicate records.
* Stores the records map in the index.

## Example

```go
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
```

On my iMac it runs pretty fast

```
Writing lorem ipsums to 'facts'

2013/09/13 14:17:21 Wrote 483 words to 'Lorem Ipsums' (took 0.000425s)
2013/09/13 14:17:21 Wrote 430 words to 'Lorem Ipsums' (took 0.000309s)
2013/09/13 14:17:21 Wrote 307 words to 'Lorem Ipsums' (took 0.000188s)
2013/09/13 14:17:21 Wrote 243 words to 'Lorem Ipsums' (took 0.000141s)

Searching for 'a needle in a Hay Stack' in 'Lorem Ipsums' index
Found 1 record(s) in 0.000001s: First Lorem Ipsum (score: 0.86)
```

## Todo
* Write tests
* Write index methods for updating and deleting records
* Turn this into a package
