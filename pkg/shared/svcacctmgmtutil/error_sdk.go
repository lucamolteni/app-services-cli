package svcacctmgmtutil

import (
	"errors"
	models "github.com/redhat-developer/app-services-cli/pkg/apisdk/svcacctmgmt/models"
)

func GetAPIErrorK(err error) *ErrorWithString {

	var redHatErrorRepresentationable models.RedHatErrorRepresentationable

	if ok := errors.As(err, &redHatErrorRepresentationable); ok {

		if redHatErrorRepresentationable.GetError() != nil {
			s := (*redHatErrorRepresentationable.GetError()).String()

			return &ErrorWithString{error_code_string: &s}

		}

		return nil

	}

	return nil

}

type ErrorWithString struct {
	error_code_string *string
}

func (e ErrorWithString) GetError() string {
	if e.error_code_string != nil {
		return *e.error_code_string
	} else {
		return ""
	}
}
