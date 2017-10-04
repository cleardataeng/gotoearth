package gotoearth

import (
	"errors"
	"fmt"
	"net/http"
)

func ExampleHTTPErr() {
	err := errors.New("oh noes")
	fmt.Println(HTTPErr(400, err))
	fmt.Println(HTTPErr(409, err))
	fmt.Println(HTTPErr(500, err))
	fmt.Println(HTTPErr(http.StatusNotFound, err))
	fmt.Println(HTTPErr(http.StatusNotImplemented, err))
	// Output:
	// [Bad Request] oh noes
	// [Conflict] oh noes
	// [Internal Server Error] oh noes
	// [Not Found] oh noes
	// [Not Implemented] oh noes
}
