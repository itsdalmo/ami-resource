package main

import (
	"encoding/json"
	"errors"
	"github.com/itsdalmo/ami-resource/models"
	"log"
	"os"
)

func main() {
	var request models.PutRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		log.Fatalf("failed to unmarshal request: %s", err)
	}
	_, err := Run(request)
	if err != nil {
		log.Fatalf("put failed: %s", err)
	}
}

// Run (business logic)
func Run(request models.PutRequest) (models.PutResponse, error) {
	var response models.PutResponse
	return response, errors.New("put is not implemented for this resource")
}
