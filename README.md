# Phrase Search

A small text searcher that is optimized to work with larger bodies of text.
It:

* Stores Documents, which are  made up of a UUID and a body of text.
* Breaks the body of text up into tokens using sentence, lowercase, and
  nGrams tokenizers. These tokens get stored into a Token Map, which
  makes sure duplicate tokens don't get added for the same document.
* Scores tokens based on how relavent they are to the Document.
* Stores the tokens into the Index, which then can be queried for
  relavent Documents.

## Example

On my iMac it runs pretty fast:

```
Writing documents to 'lorem ipsums'

2013/09/15 18:58:01 Wrote 483 words to 'Lorem Ipsums' (took 0.000361s)
2013/09/15 18:58:01 Wrote 430 words to 'Lorem Ipsums' (took 0.000321s)
2013/09/15 18:58:01 Wrote 307 words to 'Lorem Ipsums' (took 0.000207s)
2013/09/15 18:58:01 Wrote 243 words to 'Lorem Ipsums' (took 0.000149s)

Searching for 'a needle in a Hay Stack' in 'Lorem Ipsums' index
Found 1 document(s) in 0.000001s: Document{uuid:
45858f23-1e5a-11e3-81ca-901d4c8a20a4, body: "...", data: "First Lorem
Ipsum", score: 0.86}

Index now has 4 documents and 1117 tokens

Deleting found record

Delete took 0.000043s
Index now has 3 documents and 790 tokens
```

## Todo
* Write tests
* Write index methods for updating and deleting records
* Provide more settings (ie; variables for nGram algorithm)
* Turn this into a package
* Write a JSON API?
