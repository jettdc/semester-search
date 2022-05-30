package main

import (
	"bufio"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"jettdc/semester-search/gui"
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

out:
	for docIndex, d := range ds {
		sr := search.GetDocSearchResults(d, term)

	nextRes:
		for excerptIndex, result := range sr.Excerpts {
			fmt.Print("\033[H\033[2J")
			fmt.Println("\n\n\n")

			j := gui.Jumbotron{}
			j.SetXPadding(3)
			j.SetMaxWidth(120)

			j.Header = &gui.Line{Content: fmt.Sprintf("\"%s\"", sr.SearchTerm), Justify: "center"}

			j.AddLine(gui.Line{Content: fmt.Sprintf("[RESULT TYPE]: %s   [DOC #]: %d/%d   [RES #]: %d/%d", result.SearchType, docIndex+1, len(ds), excerptIndex+1, sr.NumResults), Justify: "center"})
			j.AddLine(gui.Line{Content: fmt.Sprintf("[Inspecting]: %s", d.Name), Justify: "center"})
			j.AddBlankLine()
			j.AddLine(gui.Line{Content: "Excerpt:", Justify: "left"})
			j.AddLine(gui.Line{Content: result.Content, Justify: "left"})
			j.AddBlankLine()

			j.Print()

			printCmdPrompt()
			input := bufio.NewScanner(os.Stdin)
			input.Scan()

			if input.Text() == "next" {
				break nextRes
			} else if input.Text() == "quit" || input.Text() == "q" {
				break out
			} else if input.Text() == "open" || input.Text() == "o" {
				open.Start(d.Path)
			}
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

func printCmdPrompt() {
	fmt.Println("The following commands are available:")
	fmt.Println("   -> [Enter]             : Next Result")
	fmt.Println("   -> \"next\" + [Enter]    : Next Document")
	fmt.Println("   -> \"open\" + [Enter]    : Open Current File")
	fmt.Println("   -> \"quit\" + [Enter]    : New Search\n")
	fmt.Print("CMD$ ")
}
