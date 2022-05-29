package ingest

import (
	"log"
	"os"
)

func init() {
	log.Println("Checking document ingest configurations.")

	if !fileExists("./internal") {
		log.Println("Creating internal directory")
		createInternalDirectory()
	} else {
		log.Println("Internal directory OK")
	}

	if !fileExists("./internal/ingested.json") {
		log.Println("Creating ingest cache file.")
		createIngestCache()
	} else {
		log.Println("Ingest cache OK")
	}

	if !fileExists("./documents") {
		log.Println("Creating documents directory")
		createInternalDirectory()
	} else {
		log.Println("Documents directory OK")
	}

	log.Println("Configurations OK")
}

func createInternalDirectory() {
	err := os.Mkdir("./internal", 0755)
	if err != nil {
		log.Fatal("Failed to create internal directory.")
	}
	return
}

func createDocumentsDirectory() {
	err := os.Mkdir("./documents", 0755)
	if err != nil {
		log.Fatal("Failed to create documents directory.")
	}
	return
}

func createIngestCache() {
	_, err := os.Create("./internal/ingested.json")
	if err != nil {
		log.Fatal("Failed to create ingest cache.")
	}
	return
}

func fileExists(path string) bool {
	_, e := os.Stat(path)
	if e != nil {
		if os.IsNotExist(e) {
			return false
		}
	}
	return true
}
