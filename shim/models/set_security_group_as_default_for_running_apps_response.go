// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// SetSecurityGroupAsDefaultForRunningAppsResponse set security group as default for running apps response
//
// swagger:model setSecurityGroupAsDefaultForRunningAppsResponse
type SetSecurityGroupAsDefaultForRunningAppsResponse struct {

	// The name
	Name string `json:"name,omitempty"`

	// The rules
	Rules []GenericObject `json:"rules"`

	// The running Default
	RunningDefault bool `json:"running_default,omitempty"`

	// The staging Default
	StagingDefault bool `json:"staging_default,omitempty"`
}

// Validate validates this set security group as default for running apps response
func (m *SetSecurityGroupAsDefaultForRunningAppsResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRules(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SetSecurityGroupAsDefaultForRunningAppsResponse) validateRules(formats strfmt.Registry) error {

	if swag.IsZero(m.Rules) { // not required
		return nil
	}

	for i := 0; i < len(m.Rules); i++ {

		if err := m.Rules[i].Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("rules" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *SetSecurityGroupAsDefaultForRunningAppsResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SetSecurityGroupAsDefaultForRunningAppsResponse) UnmarshalBinary(b []byte) error {
	var res SetSecurityGroupAsDefaultForRunningAppsResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}