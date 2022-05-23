package search

import (
	"jettdc/semester-search/ingest"
	"math"
	"strings"
)

type DocumentSearcher func(doc ingest.Document, term string, channel chan DocSearchResults)

type Excerpt struct {
	Content        string
	SearchPosition int
	SearchType     string
	SearchTerm     string
	SearchPriority int
}

type DocSearchResults struct {
	Excerpts       []Excerpt
	SearchTerm     string
	SearchType     string
	SearchPriority int // Lower is higher priority
	NumResults     int
}

func makeEmptyDocSearchResults(searchTerm string, searchType string, priority int) DocSearchResults {
	return DocSearchResults{[]Excerpt{}, searchTerm, searchType, priority, 0}
}

func (d *DocSearchResults) AddInheritingExcerpt(content string, searchPosition int) {
	d.Excerpts = append(
		d.Excerpts,
		Excerpt{
			content,
			searchPosition,
			d.SearchType,
			d.SearchTerm,
			d.SearchPriority},
	)
	d.NumResults += 1
}

func generateExcerptContext(doc Tokenized, search Tokenized, matchIndex int) string {
	lowerBound := int(math.Max(0.0, float64(matchIndex-ExcerptPadding)))
	upperBound := int(math.Min(float64(len(doc)-1), float64(matchIndex+len(search)+ExcerptPadding)))
	return "..." + strings.Join(doc[lowerBound:upperBound], " ") + "..."
}

// Any new excerpts that are deemed to overlap with an existing search result will be disregarded
func mergeDocSearchResults(existingResults DocSearchResults, newExcerpts []Excerpt) DocSearchResults {
	//log.Println(newExcerpts[0].MatchType)
	for _, excerpt := range newExcerpts {
		if !overlapsWithCurrentExcerpts(existingResults, excerpt) {
			existingResults.Excerpts = append(existingResults.Excerpts, excerpt)
		}
	}
	existingResults.NumResults = len(existingResults.Excerpts)
	return existingResults
}

func overlapsWithCurrentExcerpts(results DocSearchResults, excerpt Excerpt) bool {
	currentlyRecordedIndices := getExcerptIndices(results)
	excerptIndex := excerpt.SearchPosition
	for _, index := range currentlyRecordedIndices {
		if math.Abs(float64(excerptIndex-index)) <= float64(ExcerptOverlapThreshold) {
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
