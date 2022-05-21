package main

import (
	"fmt"
	"jettdc/semester-search/ingest"
	// "github.com/skratchdot/open-golang/open"
)

func main() {
	docs, err := ingest.IngestDocuments("./documents")
	if err != nil {
		fmt.Println("Broken")
	}
	
	for _, doc := range docs {
		fmt.Println(doc.Name)
		// err := open.Run(doc.Contents)
		// fmt.Println(err)
	}
}
