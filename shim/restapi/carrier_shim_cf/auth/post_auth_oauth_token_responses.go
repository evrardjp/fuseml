// Code generated by go-swagger; DO NOT EDIT.

package auth

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/suse/carrier/shim/models"
)

// PostAuthOauthTokenOKCode is the HTTP code returned for type PostAuthOauthTokenOK
const PostAuthOauthTokenOKCode int = 200

/*PostAuthOauthTokenOK successful response

swagger:response postAuthOauthTokenOK
*/
type PostAuthOauthTokenOK struct {

	/*
	  In: Body
	*/
	Payload *models.CreatesOAuthTokenResponse `json:"body,omitempty"`
}

// NewPostAuthOauthTokenOK creates PostAuthOauthTokenOK with default headers values
func NewPostAuthOauthTokenOK() *PostAuthOauthTokenOK {

	return &PostAuthOauthTokenOK{}
}

// WithPayload adds the payload to the post auth oauth token o k response
func (o *PostAuthOauthTokenOK) WithPayload(payload *models.CreatesOAuthTokenResponse) *PostAuthOauthTokenOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post auth oauth token o k response
func (o *PostAuthOauthTokenOK) SetPayload(payload *models.CreatesOAuthTokenResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostAuthOauthTokenOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}