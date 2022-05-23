package search

import (
	"jettdc/semester-search/ingest"
	"sort"
)

// When getting doc search results for all docs, do it concurrently

func GetDocSearchResults(doc ingest.Document, searchTerm string) DocSearchResults {
	// In descending order of relevance
	searchMethods := []DocumentSearcher{
		recordExactMatches,
		recordExactStemmerMatches,
		recordTrueProximityMatches,
		recordStemmerProximityMatches,
		recordLooseInstances,
	}

	results := make(chan DocSearchResults, len(searchMethods))
	for _, searchMethod := range searchMethods {
		go searchMethod(doc, searchTerm, results)
	}

	fullResults := make([]DocSearchResults, 0)
	for i := 0; i < len(searchMethods); i++ {
		fullResults = append(fullResults, <-results)
	}
	res := getCombinedSearchResults(searchTerm, fullResults)

	return res
}

func getCombinedSearchResults(searchTerm string, results []DocSearchResults) DocSearchResults {
	// Sort results by search priority
	sort.Slice(results, func(i, j int) bool {
		return results[i].SearchPriority <= results[j].SearchPriority
	})

	res := makeEmptyDocSearchResults(searchTerm, "Combined", 0)

	// mergeDocSearch is first come first serve
	for _, singleResult := range results {
		res = mergeDocSearchResults(res, singleResult.Excerpts)
	}

	return res
}

// Search for exact word matches, if searching "stronger soap", find exact matches "stronger soap"
func recordExactMatches(doc ingest.Document, term string, channel chan DocSearchResults) {
	searchResults := makeEmptyDocSearchResults(term, "Exact Match", 1)

	tokenizedDoc := getTokenizedText(doc.Contents).
		MakeLowerCase().
		Build()
	tokenizedSearch := getTokenizedText(term).
		MakeLowerCase().
		Build()

	exactMatches := getExactMatchIndices(tokenizedDoc, tokenizedSearch)

	for _, match := range exactMatches {
		searchResults.AddInheritingExcerpt(generateExcerptContext(tokenizedDoc, tokenizedSearch, match), match)
	}

	channel <- searchResults
}

// Search for stemmer matches to account for word variations
// If searching "stronger soap", find exact matches "strong soap"
func recordExactStemmerMatches(doc ingest.Document, term string, channel chan DocSearchResults) {
	searchResults := makeEmptyDocSearchResults(term, "Stemmed Exact Match", 2)

	tokenizedDoc := getTokenizedText(doc.Contents).
		MakeLowerCase().
		StemmerFilter().
		Build()
	tokenizedSearch := getTokenizedText(term).
		MakeLowerCase().
		StemmerFilter().
		Build()

	exactMatches := getExactMatchIndices(tokenizedDoc, tokenizedSearch)

	// There is a one-one mapping between the tokenized document below and stem tokenized document used for matching,
	// and so when creating our excerpt we will use the indices from the stemmed version with the
	// words from the basic version to make it readable.
	basicTokenizedDoc := getTokenizedText(doc.Contents).
		MakeLowerCase().
		Build()
	for _, match := range exactMatches {
		searchResults.AddInheritingExcerpt(generateExcerptContext(basicTokenizedDoc, tokenizedSearch, match), match)
	}

	channel <- searchResults
}

func getExactMatchIndices(tokenizedDoc Tokenized, tokenizedSearch Tokenized) []int {
	exactMatchIndices := make([]int, 0)
	for i, word := range tokenizedDoc {
		if word == tokenizedSearch[0] {
			fullMatch := true
			for j, v := range tokenizedSearch {
				if v != tokenizedDoc[i+j] { // Check boundary condition here
					fullMatch = false
				}
			}
			if fullMatch {
				exactMatchIndices = append(exactMatchIndices, i)
			}
		}
	}

	return exactMatchIndices
}

// Search for disconnected instances of exact search words
// If searching for "stronger soap", might return an excerpt with "stronger than competitors and is a leading soap"
func recordTrueProximityMatches(doc ingest.Document, term string, channel chan DocSearchResults) {
	searchResults := makeEmptyDocSearchResults(term, "Exact Proximity Match", 3)

	tokenizedDoc := getTokenizedText(doc.Contents).
		MakeLowerCase().
		Build()
	tokenizedSearch := getTokenizedText(term).
		MakeLowerCase().
		RemoveStopWords().
		Build()

	matches := getProximityMatchIndices(tokenizedDoc, tokenizedSearch)
	for _, match := range matches {
		searchResults.AddInheritingExcerpt(generateExcerptContext(tokenizedDoc, tokenizedSearch, match), match)
	}

	channel <- searchResults
}

// Search for disconnected instances of stemmer search words
// If searching for "stronger soap", might return an excerpt with "the soap is as strong as can be"
func recordStemmerProximityMatches(doc ingest.Document, term string, channel chan DocSearchResults) {
	searchResults := makeEmptyDocSearchResults(term, "Stemmed Proximity Match", 4)

	tokenizedDoc := getTokenizedText(doc.Contents).
		MakeLowerCase().
		StemmerFilter().
		Build()
	tokenizedSearch := getTokenizedText(term).
		MakeLowerCase().
		RemoveStopWords().
		StemmerFilter().
		Build()

	matches := getProximityMatchIndices(tokenizedDoc, tokenizedSearch)

	// Again, match the match indices with a non-stemmed document as to get human-readable results
	basicTokenizedDoc := getTokenizedText(doc.Contents).
		MakeLowerCase().
		Build()
	for _, match := range matches {
		searchResults.AddInheritingExcerpt(generateExcerptContext(basicTokenizedDoc, tokenizedSearch, match), match)
	}

	channel <- searchResults
}

func getProximityMatchIndices(tokenizedDoc Tokenized, tokenizedSearch Tokenized) []int {
	matches := make([]int, 0)
	for docIndex, docWord := range tokenizedDoc {
		if sliceContains(tokenizedSearch, docWord) && WordsAreInProximity(tokenizedSearch, tokenizedDoc, docIndex) {
			matches = append(matches, docIndex)
		}
	}
	return matches
}

// Pick up any missed terms, either exact or stemmer. Any instance that matches any of the search words will be returned
// If searching for "stronger soap", might return an excerpt with "the soap is amazing" or "I am a very strong guy"
// Perhaps later change so that at least half of the terms have to be found?
func recordLooseInstances(doc ingest.Document, term string, channel chan DocSearchResults) {
	searchResults := makeEmptyDocSearchResults(term, "Loose Match", 4)

	basicTokenizedDoc := getTokenizedText(doc.Contents).
		MakeLowerCase().
		Build()
	basicTokenizedSearch := getTokenizedText(term).
		MakeLowerCase().
		RemoveStopWords().
		Build()

	exactMatches := make([]int, 0)
	for i, word := range basicTokenizedDoc {
		if sliceContains(basicTokenizedSearch, word) {
			exactMatches = append(exactMatches, i)
		}
	}

	for _, match := range exactMatches {
		searchResults.AddInheritingExcerpt(generateExcerptContext(basicTokenizedDoc, basicTokenizedSearch, match), match)
	}

	channel <- searchResults
}
