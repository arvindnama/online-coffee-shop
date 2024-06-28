// Code generated by go-swagger; DO NOT EDIT.

package orders

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/arvindnama/golang-microservices/order-service/sdk/models"
)

// GetAllOrdersReader is a Reader for the GetAllOrders structure.
type GetAllOrdersReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAllOrdersReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAllOrdersOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetAllOrdersUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetAllOrdersInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /orders] getAllOrders", response, response.Code())
	}
}

// NewGetAllOrdersOK creates a GetAllOrdersOK with default headers values
func NewGetAllOrdersOK() *GetAllOrdersOK {
	return &GetAllOrdersOK{}
}

/*
GetAllOrdersOK describes a response with status code 200, with default header values.

A list of orders
*/
type GetAllOrdersOK struct {
	Payload []*models.Order
}

// IsSuccess returns true when this get all orders o k response has a 2xx status code
func (o *GetAllOrdersOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get all orders o k response has a 3xx status code
func (o *GetAllOrdersOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get all orders o k response has a 4xx status code
func (o *GetAllOrdersOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get all orders o k response has a 5xx status code
func (o *GetAllOrdersOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get all orders o k response a status code equal to that given
func (o *GetAllOrdersOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get all orders o k response
func (o *GetAllOrdersOK) Code() int {
	return 200
}

func (o *GetAllOrdersOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /orders][%d] getAllOrdersOK %s", 200, payload)
}

func (o *GetAllOrdersOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /orders][%d] getAllOrdersOK %s", 200, payload)
}

func (o *GetAllOrdersOK) GetPayload() []*models.Order {
	return o.Payload
}

func (o *GetAllOrdersOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAllOrdersUnauthorized creates a GetAllOrdersUnauthorized with default headers values
func NewGetAllOrdersUnauthorized() *GetAllOrdersUnauthorized {
	return &GetAllOrdersUnauthorized{}
}

/*
GetAllOrdersUnauthorized describes a response with status code 401, with default header values.

GetAllOrdersUnauthorized get all orders unauthorized
*/
type GetAllOrdersUnauthorized struct {
	Payload *models.ValidationError
}

// IsSuccess returns true when this get all orders unauthorized response has a 2xx status code
func (o *GetAllOrdersUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get all orders unauthorized response has a 3xx status code
func (o *GetAllOrdersUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get all orders unauthorized response has a 4xx status code
func (o *GetAllOrdersUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this get all orders unauthorized response has a 5xx status code
func (o *GetAllOrdersUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this get all orders unauthorized response a status code equal to that given
func (o *GetAllOrdersUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the get all orders unauthorized response
func (o *GetAllOrdersUnauthorized) Code() int {
	return 401
}

func (o *GetAllOrdersUnauthorized) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /orders][%d] getAllOrdersUnauthorized %s", 401, payload)
}

func (o *GetAllOrdersUnauthorized) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /orders][%d] getAllOrdersUnauthorized %s", 401, payload)
}

func (o *GetAllOrdersUnauthorized) GetPayload() *models.ValidationError {
	return o.Payload
}

func (o *GetAllOrdersUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAllOrdersInternalServerError creates a GetAllOrdersInternalServerError with default headers values
func NewGetAllOrdersInternalServerError() *GetAllOrdersInternalServerError {
	return &GetAllOrdersInternalServerError{}
}

/*
GetAllOrdersInternalServerError describes a response with status code 500, with default header values.

GetAllOrdersInternalServerError get all orders internal server error
*/
type GetAllOrdersInternalServerError struct {
	Payload *models.ValidationError
}

// IsSuccess returns true when this get all orders internal server error response has a 2xx status code
func (o *GetAllOrdersInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get all orders internal server error response has a 3xx status code
func (o *GetAllOrdersInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get all orders internal server error response has a 4xx status code
func (o *GetAllOrdersInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get all orders internal server error response has a 5xx status code
func (o *GetAllOrdersInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get all orders internal server error response a status code equal to that given
func (o *GetAllOrdersInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get all orders internal server error response
func (o *GetAllOrdersInternalServerError) Code() int {
	return 500
}

func (o *GetAllOrdersInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /orders][%d] getAllOrdersInternalServerError %s", 500, payload)
}

func (o *GetAllOrdersInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /orders][%d] getAllOrdersInternalServerError %s", 500, payload)
}

func (o *GetAllOrdersInternalServerError) GetPayload() *models.ValidationError {
	return o.Payload
}

func (o *GetAllOrdersInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
