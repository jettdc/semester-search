package main

import (
	"fmt"
	"jettdc/semester-search/ingest"
	"jettdc/semester-search/search"
	"log"
)

func main() {
	const SearchTerm = "queer cowboy"

	docs, err := ingest.IngestDocuments("./documents")
	if err != nil {
		fmt.Println("Broken")
	}

	idx := make(search.Index)
	idx.IndexDocuments(docs)
	ds := idx.Search(SearchTerm)

	for _, d := range ds {
		sr := search.GetDocSearchResults(d, SearchTerm)
		log.Println(sr.NumResults, d.Name)
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
	}
}
