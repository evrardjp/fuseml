// Code generated by go-swagger; DO NOT EDIT.

package routes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/suse/carrier/shim/models"
)

// ListAllAppsForRouteOKCode is the HTTP code returned for type ListAllAppsForRouteOK
const ListAllAppsForRouteOKCode int = 200

/*ListAllAppsForRouteOK successful response

swagger:response listAllAppsForRouteOK
*/
type ListAllAppsForRouteOK struct {

	/*
	  In: Body
	*/
	Payload *models.ListAllAppsForRouteResponsePaged `json:"body,omitempty"`
}

// NewListAllAppsForRouteOK creates ListAllAppsForRouteOK with default headers values
func NewListAllAppsForRouteOK() *ListAllAppsForRouteOK {

	return &ListAllAppsForRouteOK{}
}

// WithPayload adds the payload to the list all apps for route o k response
func (o *ListAllAppsForRouteOK) WithPayload(payload *models.ListAllAppsForRouteResponsePaged) *ListAllAppsForRouteOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list all apps for route o k response
func (o *ListAllAppsForRouteOK) SetPayload(payload *models.ListAllAppsForRouteResponsePaged) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListAllAppsForRouteOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}