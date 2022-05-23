package main

import (
	"bufio"
	"fmt"
	"jettdc/semester-search/ingest"
	"jettdc/semester-search/search"
	"log"
	"os"
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

	for docIndex, d := range ds {
		log.Println("Viewing results from", d.Name)
		sr := search.GetDocSearchResults(d, SearchTerm)
		for excerptIndex, result := range sr.Excerpts {
			log.Println("Doc ", docIndex+1, "/", len(ds), " | ", "Result", excerptIndex+1, "/", sr.NumResults)
			search.PrettyPrint(result)
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
		}
	}
}
