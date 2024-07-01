// Code generated by go-swagger; DO NOT EDIT.

package products

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewListProductsParams creates a new ListProductsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewListProductsParams() *ListProductsParams {
	return &ListProductsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewListProductsParamsWithTimeout creates a new ListProductsParams object
// with the ability to set a timeout on a request.
func NewListProductsParamsWithTimeout(timeout time.Duration) *ListProductsParams {
	return &ListProductsParams{
		timeout: timeout,
	}
}

// NewListProductsParamsWithContext creates a new ListProductsParams object
// with the ability to set a context for a request.
func NewListProductsParamsWithContext(ctx context.Context) *ListProductsParams {
	return &ListProductsParams{
		Context: ctx,
	}
}

// NewListProductsParamsWithHTTPClient creates a new ListProductsParams object
// with the ability to set a custom HTTPClient for a request.
func NewListProductsParamsWithHTTPClient(client *http.Client) *ListProductsParams {
	return &ListProductsParams{
		HTTPClient: client,
	}
}

/*
ListProductsParams contains all the parameters to send to the API endpoint

	for the list products operation.

	Typically these are written to a http.Request.
*/
type ListProductsParams struct {

	/* Currency.

	   Currency
	*/
	Currency *string

	/* PageNo.

	   page no

	   Format: int64
	   Default: 1
	*/
	PageNo *int64

	/* PageSize.

	   page size

	   Format: int64
	   Default: 10
	*/
	PageSize *int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the list products params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListProductsParams) WithDefaults() *ListProductsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the list products params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListProductsParams) SetDefaults() {
	var (
		pageNoDefault = int64(1)

		pageSizeDefault = int64(10)
	)

	val := ListProductsParams{
		PageNo:   &pageNoDefault,
		PageSize: &pageSizeDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the list products params
func (o *ListProductsParams) WithTimeout(timeout time.Duration) *ListProductsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list products params
func (o *ListProductsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list products params
func (o *ListProductsParams) WithContext(ctx context.Context) *ListProductsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list products params
func (o *ListProductsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list products params
func (o *ListProductsParams) WithHTTPClient(client *http.Client) *ListProductsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list products params
func (o *ListProductsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCurrency adds the currency to the list products params
func (o *ListProductsParams) WithCurrency(currency *string) *ListProductsParams {
	o.SetCurrency(currency)
	return o
}

// SetCurrency adds the currency to the list products params
func (o *ListProductsParams) SetCurrency(currency *string) {
	o.Currency = currency
}

// WithPageNo adds the pageNo to the list products params
func (o *ListProductsParams) WithPageNo(pageNo *int64) *ListProductsParams {
	o.SetPageNo(pageNo)
	return o
}

// SetPageNo adds the pageNo to the list products params
func (o *ListProductsParams) SetPageNo(pageNo *int64) {
	o.PageNo = pageNo
}

// WithPageSize adds the pageSize to the list products params
func (o *ListProductsParams) WithPageSize(pageSize *int64) *ListProductsParams {
	o.SetPageSize(pageSize)
	return o
}

// SetPageSize adds the pageSize to the list products params
func (o *ListProductsParams) SetPageSize(pageSize *int64) {
	o.PageSize = pageSize
}

// WriteToRequest writes these params to a swagger request
func (o *ListProductsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Currency != nil {

		// query param Currency
		var qrCurrency string

		if o.Currency != nil {
			qrCurrency = *o.Currency
		}
		qCurrency := qrCurrency
		if qCurrency != "" {

			if err := r.SetQueryParam("Currency", qCurrency); err != nil {
				return err
			}
		}
	}

	if o.PageNo != nil {

		// query param page_no
		var qrPageNo int64

		if o.PageNo != nil {
			qrPageNo = *o.PageNo
		}
		qPageNo := swag.FormatInt64(qrPageNo)
		if qPageNo != "" {

			if err := r.SetQueryParam("page_no", qPageNo); err != nil {
				return err
			}
		}
	}

	if o.PageSize != nil {

		// query param page_size
		var qrPageSize int64

		if o.PageSize != nil {
			qrPageSize = *o.PageSize
		}
		qPageSize := swag.FormatInt64(qrPageSize)
		if qPageSize != "" {

			if err := r.SetQueryParam("page_size", qPageSize); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
