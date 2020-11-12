// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// RetrievingPermissionsOnServiceInstanceResponseResource retrieving permissions on service instance response resource
//
// swagger:model retrievingPermissionsOnServiceInstanceResponseResource
type RetrievingPermissionsOnServiceInstanceResponseResource struct {

	// entity
	Entity *RetrievingPermissionsOnServiceInstanceResponse `json:"entity,omitempty"`

	// metadata
	Metadata *EntityMetadata `json:"metadata,omitempty"`
}

// Validate validates this retrieving permissions on service instance response resource
func (m *RetrievingPermissionsOnServiceInstanceResponseResource) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEntity(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMetadata(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RetrievingPermissionsOnServiceInstanceResponseResource) validateEntity(formats strfmt.Registry) error {

	if swag.IsZero(m.Entity) { // not required
		return nil
	}

	if m.Entity != nil {
		if err := m.Entity.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("entity")
			}
			return err
		}
	}

	return nil
}

func (m *RetrievingPermissionsOnServiceInstanceResponseResource) validateMetadata(formats strfmt.Registry) error {

	if swag.IsZero(m.Metadata) { // not required
		return nil
	}

	if m.Metadata != nil {
		if err := m.Metadata.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("metadata")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *RetrievingPermissionsOnServiceInstanceResponseResource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RetrievingPermissionsOnServiceInstanceResponseResource) UnmarshalBinary(b []byte) error {
	var res RetrievingPermissionsOnServiceInstanceResponseResource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}