package foo

import (
	"fmt"

	"github.com/cleardataeng/gotoearth"
)

type foo struct {
	id string
}

type Handler struct{}

func (Handler) Handle(evt interface{}) (interface{}, error) {
	var e gotoearth.Event = evt.(gotoearth.Event)
	f := foo{e.Path["fooID"]}
	fmt.Printf("handled foo: %s\n", f.id)
	return f, nil
}
