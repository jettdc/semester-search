package main

import (
	"bufio"
	"fmt"
	"jettdc/semester-search/ingest"
	"jettdc/semester-search/search"
	"os"
)

func main() {
	docs, err := ingest.IngestDocuments("./documents")
	if err != nil {
		fmt.Println("Broken")
	}

	fmt.Println("\nSemester Search is functional. Enter your search term below:")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	SearchTerm := input.Text()
	idx := make(search.Index)
	idx.IndexDocuments(docs)

	lookup(idx, SearchTerm)
}

func lookup(idx search.Index, term string) {
	ds := idx.Search(term)

	for docIndex, d := range ds {
		sr := search.GetDocSearchResults(d, term)
		for excerptIndex, result := range sr.Excerpts {
			fmt.Print("\033[H\033[2J")
			fmt.Println("Showing results for search:", sr.SearchTerm)
			fmt.Println("Doc ", docIndex+1, "/", len(ds))
			fmt.Println(result.SearchType, "\n")
			fmt.Println(d.Name, ": Result", excerptIndex+1, "/", sr.NumResults, "\n\"\"\"\n", result.Content, "\n\"\"\"")
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
		}
	}
	fmt.Print("\033[H\033[2J")
	fmt.Println("You've reached the end of your search results. Press enter to start a new search.")
	fmt.Print("[ENTER]")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Print("\033[H\033[2J")
	fmt.Println("Enter your search term below:")
	input2 := bufio.NewScanner(os.Stdin)
	input2.Scan()
	newTerm := input2.Text()
	lookup(idx, newTerm)
}
