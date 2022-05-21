package ingest

import (
	"os"
	"context"
	"io/ioutil"
	"errors"
	"log"
	"github.com/google/go-tika/tika"
)

const TIKA_PATH = "internal/tika-server-1.21.jar" 

type Document struct {
	Name string
	Path string
	Contents string
	Checksum string
}

func IngestDocuments(directory string) ([]Document, error) {
	tikaServer := setupTika()
	defer tikaServer.Stop()

	files, err := ioutil.ReadDir(directory)
    if err != nil {
		return nil, err
    }

	docs := make([]Document, len(files))
	parsedDocuments := GetParsedDocuments()

	for i, file := range files {
		path := directory + "/" + file.Name()
		checksum := GetChecksum(path)

		parsedDocument, documentHasBeenParsed := parsedDocuments[checksum]

		if !documentHasBeenParsed {
			log.Println("New document detected: ", file.Name())
			ingested := IngestDocument(tikaServer, path)
			parsedDocument = Document{file.Name(), path, ingested, checksum}
		}

		docs[i] = parsedDocument
    }

	// Dump to json before exiting
	DumpToFile(docs)

	return docs, nil
}

func DumpToFile(documents []Document) {
	json := DocumentsToJSON(documents)
	_ = ioutil.WriteFile("./internal/ingested.json", json, 0644)
}

func IngestDocument(server *tika.Server, path string) (string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	client := tika.NewClient(nil, server.URL())
	body, err := client.Parse(context.Background(), f)

	if err != nil {
		log.Fatal("Error parsing document...", err)
	}

	return body
}

func setupTika() *tika.Server {
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
        log.Println("Document parser setup failed with the following message:", a)
		os.Exit(1)
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
		return err
	}
	return nil
}



func startTikaServer() (*tika.Server, error) {
	s, err := tika.NewServer(TIKA_PATH, "")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = s.Start(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return s, nil
}