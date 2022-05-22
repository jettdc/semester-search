package search

// https://artem.krylysov.com/blog/2020/07/28/lets-build-a-full-text-search-engine/

import (
	"fmt"
	"github.com/juliangruber/go-intersect"
	"jettdc/semester-search/ingest"
	"log"
	"sort"
)

func (idx Index) Search(searchTerm string) []ingest.Document {
	var r []string
	tokenizedSearchTerm := tokenizeText(searchTerm)

	docScores := map[string]int{}

	for _, token := range tokenizedSearchTerm {
		// There are duplicates in idx[token list] so count occurrences
		ids, ok := idx[token]

		if ok {
			indvTokenDocScores := countOccurrences(ids, token)
			mergeScores(docScores, indvTokenDocScores)
			idSet := make([]string, 0)
			for key, _ := range docScores {
				idSet = append(idSet, key)
			}

			if r == nil {
				r = idSet
			} else {
				r = intersection(r, idSet)
			}
		} else {
			// Token doesn't exist.
			return nil
		}
	}

	r = sortByRelevance(r, docScores)

	log.Println(docScores)

	docs := make([]ingest.Document, len(r))
	parsedDocs := ingest.GetParsedDocuments()
	for i, doc := range r {
		parsedDoc, ok := parsedDocs[doc]
		if ok {
			docs[i] = parsedDoc
		}
	}
	return docs
}

func sortByRelevance(r []string, docScores map[string]int) []string {
	sort.SliceStable(r, func(i, j int) bool {
		return docScores[r[i]] > docScores[r[j]]
	})

	return r
}

func countOccurrences(ids []string, search string) map[string]int {
	occurrences := map[string]int{}
	for _, id := range ids {
		_, ok := occurrences[id]

		if !ok {
			occurrences[id] = 1
		} else {
			occurrences[id] += 1
		}
	}

	return occurrences
}

func mergeScores(docScores map[string]int, indvScores map[string]int) {
	for key, value := range indvScores {
		_, ok := docScores[key]

		if !ok {
			docScores[key] = value
		} else {
			docScores[key] += value
		}
	}
}

func intersection(a []string, b []string) []string {
	aInterface := make([]interface{}, len(a))
	for i, _ := range a {
		aInterface[i] = a[i]
	}

	bInterface := make([]interface{}, len(b))
	for i, _ := range b {
		bInterface[i] = b[i]
	}

	res := intersect.Hash(aInterface, bInterface)

	final := make([]string, len(res))
	for i, _ := range res {
		final[i] = fmt.Sprint(res[i])
	}

	return final
}
