package main

import (
	"encoding/json"
	"github.com/itsdalmo/ami-resource/src/models"
	"github.com/itsdalmo/ami-resource/src/out"
	"log"
	"os"
)

func main() {
	var request models.PutRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		log.Fatalf("failed to unmarshal request: %s", err)
	}
	_, err := out.Run(request)
	if err != nil {
		log.Fatalf("put failed: %s", err)
	}
}
