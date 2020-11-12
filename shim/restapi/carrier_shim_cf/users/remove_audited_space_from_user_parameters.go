// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// NewRemoveAuditedSpaceFromUserParams creates a new RemoveAuditedSpaceFromUserParams object
// no default values defined in spec.
func NewRemoveAuditedSpaceFromUserParams() RemoveAuditedSpaceFromUserParams {

	return RemoveAuditedSpaceFromUserParams{}
}

// RemoveAuditedSpaceFromUserParams contains all the bound params for the remove audited space from user operation
// typically these are obtained from a http.Request
//
// swagger:parameters removeAuditedSpaceFromUser
type RemoveAuditedSpaceFromUserParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The audited_space_guid parameter is used as a part of the request URL: '/v2/users/:guid/audited_spaces/:audited_space_guid'
	  Required: true
	  In: path
	*/
	AuditedSpaceGUID string
	/*The guid parameter is used as a part of the request URL: '/v2/users/:guid/audited_spaces/:audited_space_guid'
	  Required: true
	  In: path
	*/
	GUID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewRemoveAuditedSpaceFromUserParams() beforehand.
func (o *RemoveAuditedSpaceFromUserParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rAuditedSpaceGUID, rhkAuditedSpaceGUID, _ := route.Params.GetOK("audited_space_guid")
	if err := o.bindAuditedSpaceGUID(rAuditedSpaceGUID, rhkAuditedSpaceGUID, route.Formats); err != nil {
		res = append(res, err)
	}

	rGUID, rhkGUID, _ := route.Params.GetOK("guid")
	if err := o.bindGUID(rGUID, rhkGUID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindAuditedSpaceGUID binds and validates parameter AuditedSpaceGUID from path.
func (o *RemoveAuditedSpaceFromUserParams) bindAuditedSpaceGUID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.AuditedSpaceGUID = raw

	return nil
}

// bindGUID binds and validates parameter GUID from path.
func (o *RemoveAuditedSpaceFromUserParams) bindGUID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.GUID = raw

	return nil
}