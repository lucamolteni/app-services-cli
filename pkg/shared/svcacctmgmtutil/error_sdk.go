package svcacctmgmtutil

import (
	"errors"
	svcacctmgmtclient "github.com/redhat-developer/app-services-cli/pkg/apisdk/svcacctmgmt/models"
)

// GetAPIErrorK gets a strongly typed error from an error
func GetAPIErrorK(err error) *svcacctmgmtclient.RedHatErrorRepresentation {

	var openapiError GenericOpenAPIError

	if ok := errors.As(err, &openapiError); ok {

		rherr := svcacctmgmtclient.NewRedHatErrorRepresentation()
		getError := openapiError.GetError()

		if getError == nil {
			return nil
		}

		rherr.SetErrorDescription(getError)

		return rherr
	}

	return nil
}

// IsAPIError returns true if the error contains the errCode
func IsAPIError(err error, code string) bool {
	mappedErr := GetAPIErrorK(err)
	if mappedErr == nil {
		return false
	}

	return (*mappedErr.GetError()).String() == string(code)
}

// GenericOpenAPIError Provides access to the body, error and model on returned errors.
type GenericOpenAPIError interface {
	GetError() *string
	GetErrorDescription() *string
}
