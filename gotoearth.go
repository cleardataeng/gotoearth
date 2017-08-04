package gotoearth

import "fmt"

// Event is a downtoearth event.
type Event struct {
	Body        map[string]interface{} `json:"body"`
	Path        map[string]string      `json:"path"`
	Querystring map[string]string      `json:"querystring"`
	Route       string                 `json:"route"`
}

// Handler is used to call a method delegate based on an event.
type Handler interface {
	Handle(evt interface{}) (interface{}, error)
}

// Router is for delegating handling methods.
type Router struct {
	// Handlers are types that satisfy the Handler interface.
	// This is public so you can set it directly rather than using SetHandler.
	Handlers map[string]Handler
}

// Route the event to the correct delegate method.
func (r Router) Route(evt Event) (interface{}, error) {
	if route, ok := r.Handlers[evt.Route]; ok {
		return route.Handle(evt)
	}
	return "", fmt.Errorf("%s: no matching route", evt.Route)
}

// SetHandler adds a Handler to the Router.
// This is probably superfluous. However, there may be need for fancy controls.
func (r *Router) SetHandler(route string, handler Handler) {
	r.Handlers[route] = handler
}
