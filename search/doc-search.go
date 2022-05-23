package search

import (
	"encoding/json"
	"fmt"
	"jettdc/semester-search/ingest"
	"log"
	"math"
	"strings"
)

const PROXIMITY_RADIUS = 5

type Excerpt struct {
	Content        string
	SearchPosition int
}

type DocSearchResults struct {
	Excerpts []Excerpt
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

// When getting doc search results for all docs, do it concurrently

func GetDocSearchResults(doc ingest.Document, searchTerm string) DocSearchResults {
	res := DocSearchResults{}

	// In descending order of relevance
	res = recordExactMatches(res, doc, searchTerm)
	res = recordExactStemmerMatches(res, doc, searchTerm)
	res = recordTrueProximityMatches(res, doc, searchTerm)
	res = recordStemmerProximityMatches(res, doc, searchTerm)
	res = recordLooseInstances(res, doc, searchTerm)

	PrettyPrint(res)
	log.Println(len(res.Excerpts))

	return res
}

// Search for exact word matches, if searching "stronger soap", find exact matches "stronger soap"
func recordExactMatches(res DocSearchResults, doc ingest.Document, term string) DocSearchResults {
	basicTokenizedDoc := BasicTokenize(doc.Contents)
	basicTokenizedSearch := BasicTokenize(term)

	exactMatches := make([]int, 0)
	for i, word := range basicTokenizedDoc {
		if word == basicTokenizedSearch[0] {
			fullMatch := true
			for j, v := range basicTokenizedSearch {
				if v != basicTokenizedDoc[i+j] { // Check boundary condition here
					fullMatch = false
				}
			}
			if fullMatch {
				exactMatches = append(exactMatches, i)
			}
		}
	}
	excerpts := make([]Excerpt, 0)
	for _, match := range exactMatches {
		excerpts = append(excerpts, makeExcerpt(match, basicTokenizedDoc, basicTokenizedSearch))
	}

	return mergeDocSearchResults(res, excerpts)
}

// Search for stemmer matches to account for word variations
// If searching "stronger soap", find exact matches "strong soap"
func recordExactStemmerMatches(res DocSearchResults, doc ingest.Document, term string) DocSearchResults {
	stemTokenizedDoc := TokenizeTextKeepStopwords(doc.Contents)
	stemTokenizedSearch := TokenizeTextKeepStopwords(term)

	exactMatches := make([]int, 0)
	for i, word := range stemTokenizedDoc {
		if word == stemTokenizedSearch[0] {
			fullMatch := true
			for j, v := range stemTokenizedSearch {
				if v != stemTokenizedDoc[i+j] { // Check boundary condition here
					fullMatch = false
				}
			}
			if fullMatch {
				exactMatches = append(exactMatches, i)
			}
		}
	}

	// There is a one-one mapping between the basic tokenized document and stem tokenized document,
	// and so when creating our excerpt we will use the indices from the stemmed version with the
	// words from the basic version, to make it readable.
	basicTokenizedDoc := BasicTokenize(doc.Contents)
	excerpts := make([]Excerpt, 0)
	for _, match := range exactMatches {
		excerpts = append(excerpts, makeExcerpt(match, basicTokenizedDoc, stemTokenizedSearch))
	}

	return mergeDocSearchResults(res, excerpts)
}

// Search for disconnected instances of exact search words
// If searching for "stronger soap", might return an excerpt with "stronger than competitors and is a leading soap"
func recordTrueProximityMatches(res DocSearchResults, doc ingest.Document, term string) DocSearchResults {
	basicTokenizedDoc := BasicTokenize(doc.Contents)
	basicTokenizedSearch := BasicTokenizeNoStopwords(term)

	matches := make([]int, 0)
	for docIndex, docWord := range basicTokenizedDoc {
		if sliceContains(basicTokenizedSearch, docWord) && WordsAreInProximity(basicTokenizedSearch, basicTokenizedDoc, docIndex) {
			matches = append(matches, docIndex)
		}
	}

	excerpts := make([]Excerpt, 0)
	for _, match := range matches {
		excerpts = append(excerpts, makeExcerpt(match, basicTokenizedDoc, basicTokenizedSearch))
	}

	return mergeDocSearchResults(res, excerpts)
}

// Search for disconnected instances of stemmer search words
// If searching for "stronger soap", might return an excerpt with "the soap is as strong as can be"
func recordStemmerProximityMatches(res DocSearchResults, doc ingest.Document, term string) DocSearchResults {
	tokenizedDoc := TokenizeTextKeepStopwords(doc.Contents)
	tokenizedSearch := tokenizeText(term)

	matches := make([]int, 0)
	for docIndex, docWord := range tokenizedDoc {
		if sliceContains(tokenizedSearch, docWord) && WordsAreInProximity(tokenizedSearch, tokenizedDoc, docIndex) {
			matches = append(matches, docIndex)
		}
	}

	basicTokenizedDoc := BasicTokenize(doc.Contents)
	excerpts := make([]Excerpt, 0)
	for _, match := range matches {
		excerpts = append(excerpts, makeExcerpt(match, basicTokenizedDoc, tokenizedSearch))
	}

	return mergeDocSearchResults(res, excerpts)
}

// Pick up any missed terms, either exact or stemmer. Any instance that matches any of the search words will be returned
// If searching for "stronger soap", might return an excerpt with "the soap is amazing" or "I am a very strong guy"
func recordLooseInstances(res DocSearchResults, doc ingest.Document, term string) DocSearchResults {

	return res
}

// Return an excerpt with padding surrounding the index requested
func makeExcerpt(index int, doc []string, search []string) Excerpt {
	padding := 50
	lowerBound := int(math.Max(0.0, float64(index-padding)))
	upperBound := int(math.Min(float64(len(doc)-1), float64(index+len(search)+padding)))
	return Excerpt{"..." + strings.Join(doc[lowerBound:upperBound], " ") + "...", index}
}

// Any new excerpts that are deemed to overlap with an existing search result will be disregarded
func mergeDocSearchResults(existingResults DocSearchResults, newExcerpts []Excerpt) DocSearchResults {
	for _, excerpt := range newExcerpts {
		if !overlapsWithCurrentExcerpts(existingResults, excerpt) {
			existingResults.Excerpts = append(existingResults.Excerpts, excerpt)
		}
	}
	return existingResults
}

func overlapsWithCurrentExcerpts(results DocSearchResults, excerpt Excerpt) bool {
	overlapThreshold := 20
	currentlyRecordedIndices := getExcerptIndices(results)
	excerptIndex := excerpt.SearchPosition
	for _, index := range currentlyRecordedIndices {
		if math.Abs(float64(excerptIndex-index)) <= float64(overlapThreshold) {
			return true
		}
	}
	return false
}

func getExcerptIndices(results DocSearchResults) []int {
	indices := make([]int, 0)
	for _, excerpt := range results.Excerpts {
		indices = append(indices, excerpt.SearchPosition)
	}
	return indices
}
