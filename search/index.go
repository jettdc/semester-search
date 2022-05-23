package search

import (
	"jettdc/semester-search/ingest"
)

type Index map[string][]string

func (idx Index) IndexDocuments(docs []ingest.Document) {
	for _, doc := range docs {
		idx.indexDocument(doc)
	}
}

func (idx Index) indexDocument(doc ingest.Document) {
	tokenizedDocument := getTokenizedText(doc.Contents).
		MakeLowerCase().
		RemoveStopWords().
		StemmerFilter().
		Build()
	for _, token := range tokenizedDocument {
		ids := idx[token]
		idx[token] = append(ids, doc.Checksum)
	}
}
