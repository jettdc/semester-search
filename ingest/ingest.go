package ingest

import (
	"context"
	"errors"
	"github.com/google/go-tika/tika"
	"io/ioutil"
	"log"
	"os"
)

const TIKA_PATH = "internal/tika-server-1.21.jar"

func IngestDocuments(directory string) ([]Document, error) {
	log.Println("Ingesting documents")

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	docs := make([]Document, len(files))
	parsedDocuments := GetParsedDocuments()

	var tikaServer *tika.Server

	for i, file := range files {
		path := directory + "/" + file.Name()
		checksum := GetChecksum(path)

		parsedDocument, documentHasBeenParsed := parsedDocuments[checksum]

		if !documentHasBeenParsed {
			log.Println("New document detected: ", file.Name())

			// Only start the document parsing server if needed
			if tikaServer == nil {
				tikaServer = setupTika()
				defer tikaServer.Stop()
			}

			ingested := IngestDocument(tikaServer, path)
			parsedDocument = Document{file.Name(), path, ingested, checksum}
		}

		docs[i] = parsedDocument
	}

	if tikaServer == nil && len(docs) > 0 {
		log.Println("No new documents detected.")
	} else {
		if len(docs) == 0 {
			log.Println("No documents found. Try placing the files into the ./documents directory.")
		}

		log.Println("Updating document cache.")
		DumpToFile(docs)
	}

	return docs, nil
}

func IngestDocument(server *tika.Server, path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.Println("Attempting to parse document located at", path)
	client := tika.NewClient(nil, server.URL())
	body, err := client.Parse(context.Background(), f)

	if err != nil {
		log.Fatal("Error parsing document...", err)
	}

	return body
}

func setupTika() *tika.Server {
	log.Println("Setting up the document parsing server")
	defer handleSetupFailure()

	// Download the document parsing server
	if !tikaDownloaded() {
		log.Println("Downloading the document parser...")
		err := downloadTika()
		if err != nil {
			panic("Something went wrong downloading the document parser.")
		}

	}

	// Run and get the document parsing server
	log.Println("Starting the document parsing server...")
	s, err := startTikaServer()
	if err != nil {
		panic("Could not start the document parsing server.")
	}

	return s
}

func handleSetupFailure() {
	if a := recover(); a != nil {
		log.Fatal("Document parser setup failed with the following message:", a)
	}
}

func tikaDownloaded() bool {
	_, err := os.Stat(TIKA_PATH)
	return !errors.Is(err, os.ErrNotExist)
}

func downloadTika() error {
	err := tika.DownloadServer(context.Background(), "1.21", TIKA_PATH)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func startTikaServer() (*tika.Server, error) {
	s, err := tika.NewServer(TIKA_PATH, "")
	if err != nil {
		log.Fatal(err)
	}

	err = s.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return s, nil
}
