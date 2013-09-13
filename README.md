# Phrase Search

A small text searcher that is optimized to work with larger bodies of text.
It:

* Stores Documents, which are  made up of a UUID and a body of text.
* Breaks the body of text up into tokens using sentence, lowercase, and
  nGrams tokenizers. These tokens get stored into a Token Map, which
  makes sure duplicate tokens don't get added for the same document.
* Tokens reference Documents.
* Stores the tokens into the Index.

## Example

On my iMac it runs pretty fast:

```
Writing lorem ipsums to 'facts'

2013/09/13 17:27:46 Wrote 483 words to 'Lorem Ipsums' (took 0.000343s)
2013/09/13 17:27:46 Wrote 430 words to 'Lorem Ipsums' (took 0.000318s)
2013/09/13 17:27:46 Wrote 307 words to 'Lorem Ipsums' (took 0.000207s)
2013/09/13 17:27:46 Wrote 243 words to 'Lorem Ipsums' (took 0.000136s)

Searching for 'a needle in a Hay Stack' in 'Lorem Ipsums' index
Found 1 document(s) in 0.000001s: Document{uuid: 54ff2ee3-1cbb-11e3-8bb7-fa8884c8312b, data: "First Lorem Ipsum", score: 0.86}
```

## Todo
* Write tests
* Write index methods for updating and deleting records
* Provide more settings (ie; variables for nGram algorithm)
* Turn this into a package
* Write a JSON API?
