// Code generated by go-swagger; DO NOT EDIT.

package pet

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// DeletePetBadRequestCode is the HTTP code returned for type DeletePetBadRequest
const DeletePetBadRequestCode int = 400

/*DeletePetBadRequest Invalid pet value

swagger:response deletePetBadRequest
*/
type DeletePetBadRequest struct {
}

// NewDeletePetBadRequest creates DeletePetBadRequest with default headers values
func NewDeletePetBadRequest() *DeletePetBadRequest {

	return &DeletePetBadRequest{}
}

// WriteResponse to the client
func (o *DeletePetBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}
