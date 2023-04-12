package kapi

import (
	kafkamgmtapi "github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt/api"
	svcacctmgmtapis "github.com/redhat-developer/app-services-cli/pkg/apisdk/svcacctmgmt/apis"
)

type KiotaAPI interface {
	KafkaMgmt() *kafkamgmtapi.Kafkas_mgmtRequestBuilder
	ServiceAccountMgmt() *svcacctmgmtapis.Service_accountsRequestBuilder
}
