package main

import (
	"fmt"
	"jettdc/semester-search/ingest"
	"jettdc/semester-search/search"
	"log"
	// "github.com/skratchdot/open-golang/open"
)

func main() {
	docs, err := ingest.IngestDocuments("./documents")
	if err != nil {
		fmt.Println("Broken")
	}

	idx := make(search.Index)
	idx.IndexDocuments(docs)
	ds := idx.Search("narratorial function")
	for i, d := range ds {
		log.Println(i, ":", d.Checksum, d.Name)
	}
}
