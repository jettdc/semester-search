package ingest

import (
	"log"
	"encoding/json"
	"encoding/hex"
	"io"
	"io/ioutil"

	"crypto/sha256"
	"os"
)


func DocumentsToJSON(documents []Document) []byte {
	parsedDocuments := make(map[string]Document)

	for _, doc := range documents {
		parsedDocuments[doc.Checksum] = doc
	}

	res, err := json.Marshal(parsedDocuments)

	if err != nil {
        log.Fatal(err)
    }

	return res
}

func GetChecksum(path string) (string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	checksum := hex.EncodeToString(h.Sum(nil))
	return checksum
}

func GetParsedDocuments() map[string]Document {
	content, err := ioutil.ReadFile("./internal/ingested.json")
    if err != nil {
        log.Fatal("Error when opening file: ", err)
    }
 
    var payload map[string]Document
    err = json.Unmarshal(content, &payload)
    if err != nil {
        log.Fatal("Error during Unmarshal(): ", err)
    }

	return payload
}

func DumpToFile(documents []Document) {
	json := DocumentsToJSON(documents)
	_ = ioutil.WriteFile("./internal/ingested.json", json, 0644)
}