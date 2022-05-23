package main

import (
	"jettdc/semester-search/search"
	"log"
)

func main() {
	doc := []string{"the", "cats", "and", "I", "are", "in", "love"}
	log.Println(search.WordsAreInProximity([]string{"I", "love", "cats"}, doc, 1))
	//docs, err := ingest.IngestDocuments("./documents")
	//if err != nil {
	//	fmt.Println("Broken")
	//}
	//
	//idx := make(search.Index)
	//idx.IndexDocuments(docs)
	//ds := idx.Search("queer cowboy")
	//
	//for i, d := range ds {
	//	log.Println(i, ":", d.Checksum, d.Name)
	//}
	//
	//search.GetDocSearchResults(ds[0], "Toxic Masculinity")

	// This call would be made after the for loop finds the word cats in the doc
}
