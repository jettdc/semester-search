package search

import (
	snowballeng "github.com/kljensen/snowball/english"
	"jettdc/semester-search/ingest"
	"strings"
	"unicode"
)

type Index map[string][]string

func (idx Index) IndexDocuments(docs []ingest.Document) {
	for _, doc := range docs {
		idx.indexDocument(doc)
	}
}

func (idx Index) indexDocument(doc ingest.Document) {
	for _, token := range tokenizeDocument(doc) {
		ids := idx[token]
		idx[token] = append(ids, doc.Checksum)
	}
}

func tokenizeDocument(document ingest.Document) []string {
	return tokenizeText(document.Contents)
}

func tokenizeText(content string) []string {
	tokens := tokenize(content)
	tokens = lowercaseFilter(tokens)
	tokens = stopwordFilter(tokens)
	tokens = stemmerFilter(tokens)
	return tokens
}

func TokenizeTextKeepStopwords(content string) []string {
	tokens := tokenize(content)
	tokens = lowercaseFilter(tokens)
	tokens = stemmerFilter(tokens)
	return tokens
}

func BasicTokenize(content string) []string {
	tokens := tokenize(content)
	tokens = lowercaseFilter(tokens)
	return tokens
}

// Split on any character that is not a letter or a number
func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}

var stopwords = map[string]struct{}{
	"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
	"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
}

func stopwordFilter(tokens []string) []string {
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := stopwords[token]; !ok {
			r = append(r, token)
		}
	}
	return r
}

func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}
	return r
}
