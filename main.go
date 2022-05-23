package main

import (
	"fmt"
	"jettdc/semester-search/ingest"
	"jettdc/semester-search/search"
	"log"
)

func main() {
	docs, err := ingest.IngestDocuments("./documents")
	if err != nil {
		fmt.Println("Broken")
	}

	idx := make(search.Index)
	idx.IndexDocuments(docs)
	ds := idx.Search("queer cowboy")

	for i, d := range ds {
		log.Println(i, ":", d.Checksum, d.Name)
	}

	search.GetDocSearchResults(ds[0], "Toxic Masculinity")
}
