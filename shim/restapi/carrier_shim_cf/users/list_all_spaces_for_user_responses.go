// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/suse/carrier/shim/models"
)

// ListAllSpacesForUserOKCode is the HTTP code returned for type ListAllSpacesForUserOK
const ListAllSpacesForUserOKCode int = 200

/*ListAllSpacesForUserOK successful response

swagger:response listAllSpacesForUserOK
*/
type ListAllSpacesForUserOK struct {

	/*
	  In: Body
	*/
	Payload *models.ListAllSpacesForUserResponsePaged `json:"body,omitempty"`
}

// NewListAllSpacesForUserOK creates ListAllSpacesForUserOK with default headers values
func NewListAllSpacesForUserOK() *ListAllSpacesForUserOK {

	return &ListAllSpacesForUserOK{}
}

// WithPayload adds the payload to the list all spaces for user o k response
func (o *ListAllSpacesForUserOK) WithPayload(payload *models.ListAllSpacesForUserResponsePaged) *ListAllSpacesForUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list all spaces for user o k response
func (o *ListAllSpacesForUserOK) SetPayload(payload *models.ListAllSpacesForUserResponsePaged) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListAllSpacesForUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}