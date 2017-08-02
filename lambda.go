package gotoearth

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// LambdaHandler is a helper for event that need only to invoke another Lambda.
type LambdaHandler struct {
	// This is the InvokeInput object documented here:
	// https://docs.aws.amazon.com/sdk-for-go/api/service/lambda/#InvokeInput
	Input lambda.InvokeInput
}

func getSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}

// Handle is the method which causes LambdaHandler to satisfy the Handler
// interface. Be sure to set the lambda.InvokeInput on LambdaHandler. There is
// likely no reason to provide the Payload in lambda.InvokeInput. If you do, it
// will be used. If you do not, the gotoearth.Event will be passed along.
func (h LambdaHandler) Handle(evt interface{}) (interface{}, error) {
	if h.Input.FunctionName == nil {
		return "", errors.New("no lambda.InvokeInput.FunctionName given")
	}
	svc := lambda.New(getSession())
	if h.Input.Payload == nil {
		payload, err := json.Marshal(evt)
		if err != nil {
			return "", fmt.Errorf("evt failed to marshal: %s", err.Error())
		}
		h.Input.Payload = payload
	}
	result, err := svc.Invoke(&h.Input)
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
