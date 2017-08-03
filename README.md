gotoearth
=========

gotoearth is a Go library for working with [downtoearth](https://github.com/cleardataeng/downtoearth) API's. You would use this in the AWS Lambda to which the downtoearth.json points your API Gateway requests. That AWS Lambda would be running the [aws-lambda-go-shim](https://github.com/eawsy/aws-lambda-go-shim), which is awesome by the way, and you should be using it for running Go in AWS Lambda.

It is named as is, rather that downtoearth-go or something similar, so that the import does not require an alias.

Installation
------------

Installation is simple:

    go get github.com/cleardataeng/gotoearth

Usage
-----

Using gotoearth is fairly simple, too. There is an example provided in the repository from which you can take guidance.

### Running code in the same invocation (local) ###

The general idea is to take a downtoearth event, and route it to delegate methods. For example, somebody makes a `GET` request at `/foo/{fooID}`, and you would like to route this to handler in the same project. To do you so, you would need the following in your root handler.

``` go
func Handle(evt gotoearth.Event, ctx *runtime.Context) (interface{}, error) {
	r := gotoearth.Router{Handlers: map[string]gotoearth.Handler{
		"GET:/foo/{fooID}": foo.Handler{},
	}}
	return r.Route(evt)
}
```

and in the `foo` package, something like this:

``` go
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
```

### Invoke another Lambda ###

Alternatively, you could have the root handler delegate by invoking another Lambda with the same event. Let's assume you have another `GET` request at `/bar/{barID}`, and you want it to invoke another AWS Lambda. In that case, you could do the following.

``` go
func Handle(evt gotoearth.Event, ctx *runtime.Context) (interface{}, error) {
	r := gotoearth.Router{Handlers: map[string]gotoearth.Handler{
		"GET:/bar/{barID}": gotoearth.LambdaHandler{lambda.InvokeInput{
			FunctionName:   aws.String("arn:aws:lambda:us-west-2:1234567890:function:bar"),
			InvocationType: aws.String("Event"),
		}},
	}}
	return r.Route(evt)
}
```

Then, in the invoked AWS Lambda, you just gotoearth.Event as the type for the event again.

_Even better yet_: gotoearth also provides the type SimpleLambda. This makes a few assumptions which make usage even easier. Since invoking a Lambda as Request / Response would be similar to just making the library call in this code, gotoearth assumes that if you are invoking another Lambda, you probably want the InvocationType to be "Event". We also just assume you are passing along the same payload and accepting all of lambda.InvokeInput defaults. Therefore, you can just give the initialization of the type a value (string not *string) for FunctionName, which could be just the function name of or a full ARN.

``` go
func Handle(evt gotoearth.Event, ctx *runtime.Context) (interface{}, error) {
	r := gotoearth.Router{Handlers: map[string]gotoearth.Handler{
		"GET:/baz/{bazID}": gotoearth.SimpleLambda{"baz"},
	}}
	return r.Route(evt)
```

How neat is that?
