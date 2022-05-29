# Semester Search
Semester search is a utility for quickly searching through downloadable class materials so that you can spend more time
learning and less time clicking through dozens of links on your professor's websites.

## How To Run
- Make sure java is installed on your computer and available in your path
- Create a directory at `./documents` to place your documents in
  - Conversely, if you run the program once without documents, it will create this folder automatically.
- run `main.go`

## How does it work?
After placing your documents (tested with pdf, pptx, and doc so far) into the `documents` directory,
the program will scan through them, noting any new documents since the last time you opened the search utility.

If you have not added any new documents, the engine will use the cached versions of the (parsed)
documents to perform your searches.

If you have added or removed documents since your last search, the engine will reindex your documents. It does this by
starting a Tika server (including downloading it if you do not have it installed), and then feeding each of your
documents to the service. The server responds with the body of the document, which is then stored for you to search.

After loading the documents into memory (either from cache or via parsing), the engine will create a full text search
index from their content. Using this index, users can search to find documents that contain what they are looking for.
Results at this point are sorted by hits per document.

To dive deeper, a further text search is performed on each document to get specific excerpts that you can peruse to make
sure that you are looking at the correct document. These document specific search methods include exact phrase matching,
stemmed phrase matching, search term proximity matching, and loose term matching.

When you've found the correct document (and viewing just the excerpt
is insufficient), you can use the dedicated keyboard shortcut to open it in your favorite document viewer.


*Note: this project requires Java to run, as the Tika document parsing server depends on it*