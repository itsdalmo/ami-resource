package out

import (
	"errors"
	"github.com/itsdalmo/ami-resource/src/models"
)

// Run (business logic)
func Run(request models.PutRequest) (models.PutResponse, error) {
	var response models.PutResponse
	return response, errors.New("put is not implemented for this resource")
}
