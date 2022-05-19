package main

import (
	"fmt"
	"jettdc/semester-search/ingest"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	docs, err := ingest.IngestDocuments("./documents")
	if err != nil {
		fmt.Println("Broken")
	}
	
	for _, doc := range docs {
		fmt.Println(doc.Path)
		err := open.Run(doc.Path)
		fmt.Println(err)
	}
}
