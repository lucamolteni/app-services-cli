package kafkautil

import (
	"errors"
	"github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt/models"
)

func GetAPIError(err error) *ErrorWithCode {

	var kafkaError models.Errorable

	if ok := errors.As(err, &kafkaError); ok {

		if kafkaError.GetCode() != nil {
			s := *kafkaError.GetCode()

			return &ErrorWithCode{errorCodeString: &s}

		}

		return nil

	}

	return nil

}

type ErrorWithCode struct {
	errorCodeString *string
}

func (e ErrorWithCode) GetCode() string {
	if e.errorCodeString != nil {
		return *e.errorCodeString
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
