package kafkautil

import (
	"errors"
	"github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt/models"
)

func GetAPIError(err error) *ErrorWithCode {

	var kafkaError models.Errorable

	if ok := errors.As(err, &kafkaError); ok {
		return &ErrorWithCode{errorCodeString: kafkaError.GetCode(), reason: kafkaError.GetReason()}
	}

	return nil

}

type ErrorWithCode struct {
	errorCodeString *string
	reason          *string
}

func (e ErrorWithCode) GetCode() string {
	if e.errorCodeString != nil {
		return *e.errorCodeString
	} else {
		return ""
	}
}

func (e ErrorWithCode) GetReason() string {
	if e.reason != nil {
		return *e.reason
	} else {
		return ""
	}
}

// IsAPIError returns true if the error contains the errCode
func IsAPIError(err error, code string) bool {
	mappedErr := GetAPIError(err)
	if mappedErr == nil {
		return false
	}

	return mappedErr.GetCode() == string(code)
}
