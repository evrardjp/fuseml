// Code generated by go-swagger; DO NOT EDIT.

package stacks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// RetrieveStackHandlerFunc turns a function with the right signature into a retrieve stack handler
type RetrieveStackHandlerFunc func(RetrieveStackParams) middleware.Responder

// Handle executing the request and returning a response
func (fn RetrieveStackHandlerFunc) Handle(params RetrieveStackParams) middleware.Responder {
	return fn(params)
}

// RetrieveStackHandler interface for that can handle valid retrieve stack params
type RetrieveStackHandler interface {
	Handle(RetrieveStackParams) middleware.Responder
}

// NewRetrieveStack creates a new http.Handler for the retrieve stack operation
func NewRetrieveStack(ctx *middleware.Context, handler RetrieveStackHandler) *RetrieveStack {
	return &RetrieveStack{Context: ctx, Handler: handler}
}

/*RetrieveStack swagger:route GET /stacks/{guid} stacks retrieveStack

Retrieve a Particular Stack

curl --insecure -i %s/v2/stacks/{guid} -X GET -H 'Authorization: %s'

*/
type RetrieveStack struct {
	Context *middleware.Context
	Handler RetrieveStackHandler
}

func (o *RetrieveStack) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewRetrieveStackParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}