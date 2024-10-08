// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ProductsPaginatedResponse ProductsPaginatedResponse products paginated response
//
// swagger:model ProductsPaginatedResponse
type ProductsPaginatedResponse struct {

	// has more
	HasMore bool `json:"hasMore,omitempty"`

	// page no
	PageNo int64 `json:"pageNo,omitempty"`

	// page size
	PageSize int64 `json:"pageSize,omitempty"`

	// content
	Content Products `json:"content,omitempty"`
}

// Validate validates this products paginated response
func (m *ProductsPaginatedResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateContent(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProductsPaginatedResponse) validateContent(formats strfmt.Registry) error {
	if swag.IsZero(m.Content) { // not required
		return nil
	}

	if err := m.Content.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("content")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("content")
		}
		return err
	}

	return nil
}

// ContextValidate validate this products paginated response based on the context it is used
func (m *ProductsPaginatedResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateContent(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProductsPaginatedResponse) contextValidateContent(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Content.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("content")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("content")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProductsPaginatedResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProductsPaginatedResponse) UnmarshalBinary(b []byte) error {
	var res ProductsPaginatedResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
