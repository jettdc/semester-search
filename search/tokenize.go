package search

import (
	snowballeng "github.com/kljensen/snowball/english"
	"strings"
	"unicode"
)

type Tokenized []string

func getTokenizedText(text string) *Tokenized {
	tokenizedText := Tokenized(strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	}))

	return &tokenizedText
}

func (t *Tokenized) Build() Tokenized {
	return *t
}

func (t *Tokenized) MakeLowerCase() *Tokenized {
	r := make([]string, len(*t))
	for i, token := range *t {
		r[i] = strings.ToLower(token)
	}
	*t = r
	return t
}

func (t *Tokenized) StemmerFilter() *Tokenized {
	r := make([]string, len(*t))
	for i, token := range *t {
		r[i] = snowballeng.Stem(token, false)
	}
	*t = r
	return t
}

func (t *Tokenized) RemoveStopWords() *Tokenized {
	r := make([]string, 0, len(*t))
	for _, token := range *t {
		if _, ok := stopwords[token]; !ok {
			r = append(r, token)
		}
	}
	*t = r
	return t
}

var stopwords = map[string]struct{}{
	"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
	"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
}
