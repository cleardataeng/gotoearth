package gotoearth

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// Lambda is a helper for event that need only to invoke another Lambda.
type Lambda struct {
	// This is the InvokeInput object documented here:
	// https://docs.aws.amazon.com/sdk-for-go/api/service/lambda/#InvokeInput
	Input lambda.InvokeInput
}

// SimpleLambda is similar to Lamba but requires only an FunctionName.
// It assumes that the InvocationType will be "Event". It will use the passalong
// payload of the current Event. It will use default values for all other values
// in the lambda.InvokeInput. In addition, you need only to pass in a string.
// This will be converted to aws.String for you.
type SimpleLambda struct {
	// This is the name of the function or the full ARN. It is a direct
	// passthrough to the lambda.InvokeInput.FunctionName. You can view that
	// documentation at: https://docs.aws.amazon.com/sdk-for-go/api/service/lambda/#InvokeInput.
	// The only different is that this field is a string, not *string.
	FunctionName string
}

func getSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}

func invoke(input lambda.InvokeInput) (interface{}, error) {
	svc := lambda.New(getSession())
	result, err := svc.Invoke(&input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lambda.ErrCodeServiceException:
				return "", fmt.Errorf(lambda.ErrCodeServiceException, aerr.Error())
			case lambda.ErrCodeResourceNotFoundException:
				return "", fmt.Errorf(lambda.ErrCodeResourceNotFoundException, aerr.Error())
			case lambda.ErrCodeInvalidRequestContentException:
				return "", fmt.Errorf(lambda.ErrCodeInvalidRequestContentException, aerr.Error())
			case lambda.ErrCodeRequestTooLargeException:
				return "", fmt.Errorf(lambda.ErrCodeRequestTooLargeException, aerr.Error())
			case lambda.ErrCodeUnsupportedMediaTypeException:
				return "", fmt.Errorf(lambda.ErrCodeUnsupportedMediaTypeException, aerr.Error())
			case lambda.ErrCodeTooManyRequestsException:
				return "", fmt.Errorf(lambda.ErrCodeTooManyRequestsException, aerr.Error())
			case lambda.ErrCodeInvalidParameterValueException:
				return "", fmt.Errorf(lambda.ErrCodeInvalidParameterValueException, aerr.Error())
			case lambda.ErrCodeEC2UnexpectedException:
				return "", fmt.Errorf(lambda.ErrCodeEC2UnexpectedException, aerr.Error())
			case lambda.ErrCodeSubnetIPAddressLimitReachedException:
				return "", fmt.Errorf(lambda.ErrCodeSubnetIPAddressLimitReachedException, aerr.Error())
			case lambda.ErrCodeENILimitReachedException:
				return "", fmt.Errorf(lambda.ErrCodeENILimitReachedException, aerr.Error())
			case lambda.ErrCodeEC2ThrottledException:
				return "", fmt.Errorf(lambda.ErrCodeEC2ThrottledException, aerr.Error())
			case lambda.ErrCodeEC2AccessDeniedException:
				return "", fmt.Errorf(lambda.ErrCodeEC2AccessDeniedException, aerr.Error())
			case lambda.ErrCodeInvalidSubnetIDException:
				return "", fmt.Errorf(lambda.ErrCodeInvalidSubnetIDException, aerr.Error())
			case lambda.ErrCodeInvalidSecurityGroupIDException:
				return "", fmt.Errorf(lambda.ErrCodeInvalidSecurityGroupIDException, aerr.Error())
			case lambda.ErrCodeInvalidZipFileException:
				return "", fmt.Errorf(lambda.ErrCodeInvalidZipFileException, aerr.Error())
			case lambda.ErrCodeKMSDisabledException:
				return "", fmt.Errorf(lambda.ErrCodeKMSDisabledException, aerr.Error())
			case lambda.ErrCodeKMSInvalidStateException:
				return "", fmt.Errorf(lambda.ErrCodeKMSInvalidStateException, aerr.Error())
			case lambda.ErrCodeKMSAccessDeniedException:
				return "", fmt.Errorf(lambda.ErrCodeKMSAccessDeniedException, aerr.Error())
			case lambda.ErrCodeKMSNotFoundException:
				return "", fmt.Errorf(lambda.ErrCodeKMSNotFoundException, aerr.Error())
			case lambda.ErrCodeInvalidRuntimeException:
				return "", fmt.Errorf(lambda.ErrCodeInvalidRuntimeException, aerr.Error())
			default:
				return "", aerr
			}
		}
		return "", err
	}
	return result, nil
}

// Handle is the method which causes Lambda to satisfy the Handler interface.
// Be sure to set Input (lambda.InvokeInput). There is likely no reason to
// provide the Payload in Input. If you do, it will be used. If you do not, the
// gotoearth.Event will be passed along.
func (l Lambda) Handle(evt interface{}) (interface{}, error) {
	if l.Input.FunctionName == nil {
		return "", errors.New("no lambda.InvokeInput.FunctionName given")
	}
	if l.Input.Payload == nil {
		payload, err := json.Marshal(evt)
		if err != nil {
			return "", fmt.Errorf("evt failed to marshal: %s", err.Error())
		}
		l.Input.Payload = payload
	}
	return invoke(l.Input)
}

// Handle is the method which causes SimpleHandler to satisfy the Handler
// interface. The only difference is this makes assumptions and therefore
// makes usage easier.
func (l SimpleLambda) Handle(evt interface{}) (interface{}, error) {
	if l.FunctionName == "" {
		return "", errors.New("no FunctionName given")
	}
	payload, err := json.Marshal(evt)
	if err != nil {
		return "", fmt.Errorf("evt failed to marshal: %s", err.Error())
	}
	input := lambda.InvokeInput{
		FunctionName:   aws.String(l.FunctionName),
		InvocationType: aws.String("Event"),
		Payload:        payload,
	}
	return invoke(input)
}
