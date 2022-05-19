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
}

func IngestDocuments(directory string) ([]Document, error) {
	tikaServer := setupTika()
	defer tikaServer.Stop()

	files, err := ioutil.ReadDir(directory)
    if err != nil {
		return nil, err
    }

	docs := make([]Document, len(files))

	for i, file := range files {
        log.Println("Parsing ", file.Name())
		res := IngestDocument(tikaServer, directory + "/" + file.Name())
		docs[i] = Document{file.Name(), directory + "/" + file.Name(), res}
    }

	return docs, nil
}

func IngestDocument(server *tika.Server, path string) (string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return "Fixk!"
	}
	defer f.Close()

	client := tika.NewClient(nil, server.URL())
	body, err := client.Parse(context.Background(), f)

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