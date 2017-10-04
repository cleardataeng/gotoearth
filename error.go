package gotoearth

import (
	"fmt"
	"net/http"
)

// HTTPErr returns an error with the text updated as required by downtoearth.
//
// For example, passing in a 404 and an error with message "it isn't there"
// yields a new error with message "[Not Found] it isn't there". This is how
// downtoearth configures API Gateway to return the correct error.
//
// You can also pass in the constants defined in the net/http package that begin
// with "Status". For example, passing in http.StatusNotFound in the first
// argument would yield the same result as above.
//
// You can find the constants defined in the net/http package here:
// https://golang.org/pkg/net/http/#pkg-constants
func HTTPErr(code int, err error) error {
	if txt := http.StatusText(code); txt != "" {
		return fmt.Errorf("[%s] %s", txt, err)
	}
	return err
}
