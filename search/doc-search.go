package search

import (
	"jettdc/semester-search/ingest"
	"log"
	"math"
	"strings"
)

type Excerpt struct {
	Content             string
	SearchStartPosition int
}

type DocSearchResults struct {
	Excerpts []Excerpt
}

func GetDocSearchResults(doc ingest.Document, searchTerm string) DocSearchResults {
	res := DocSearchResults{}
	res = recordExactMatches(res, doc, searchTerm)
	res = recordExactStemmerMatches(res, doc, searchTerm)
	res = recordTrueProximityMatches(res, doc, searchTerm)
	res = recordStemmerProximityMatches(res, doc, searchTerm)
	res = recordLooseInstances(res, doc, searchTerm)
	log.Println(res)
	return res
}

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

func recordExactStemmerMatches(res DocSearchResults, doc ingest.Document, term string) DocSearchResults {
	return res
}

func recordTrueProximityMatches(res DocSearchResults, doc ingest.Document, term string) DocSearchResults {

	return res
}

func recordStemmerProximityMatches(res DocSearchResults, doc ingest.Document, term string) DocSearchResults {

	return res
}

func recordLooseInstances(res DocSearchResults, doc ingest.Document, term string) DocSearchResults {

	return res
}

func makeExcerpt(index int, doc []string, search []string) Excerpt {
	padding := 50
	lowerBound := int(math.Max(0.0, float64(index-padding)))
	upperBound := int(math.Min(float64(len(doc)-1), float64(index+len(search)+padding)))
	return Excerpt{strings.Join(doc[lowerBound:upperBound], " "), index}
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
	excerptIndex := excerpt.SearchStartPosition
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
		indices = append(indices, excerpt.SearchStartPosition)
	}
	return indices
}
