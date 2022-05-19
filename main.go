package main

import (
	"fmt"
	"jettdc/semester-search/ingest"
)

func main() {
	docs, err := ingest.IngestDocuments("./documents")
	if err != nil {
		fmt.Println("Broken")
	}
	
	for _, doc := range docs {
		fmt.Println(doc.Name)
	}
}
