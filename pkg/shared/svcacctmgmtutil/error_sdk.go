package svcacctmgmtutil

import (
	"errors"
	"github.com/redhat-developer/app-services-cli/pkg/apisdk/svcacctmgmt/models"
)

func GetAPIError(err error) *ErrorWithString {

	var redHatErrorRepresentationable models.RedHatErrorRepresentationable

	if ok := errors.As(err, &redHatErrorRepresentationable); ok {

		if redHatErrorRepresentationable.GetError() != nil {
			s := (*redHatErrorRepresentationable.GetError()).String()

			return &ErrorWithString{errorCodeString: &s}

		}

		return nil

	}

	return nil

}

type ErrorWithString struct {
	errorCodeString *string
}

func (e ErrorWithString) GetError() string {
	if e.errorCodeString != nil {
		return *e.errorCodeString
	} else {
		return ""
	}
}
