package kapi

import (
	kafkamgmtapi "github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt/api"
	svcacctmgmtapis "github.com/redhat-developer/app-services-cli/pkg/apisdk/svcacctmgmt/apis"
	"github.com/redhat-developer/app-services-cli/pkg/shared/connection/api"
)

type KiotaAPI interface {
	KafkaMgmt() *kafkamgmtapi.Kafkas_mgmtRequestBuilder
	ServiceAccountMgmt() *svcacctmgmtapis.Service_accountsRequestBuilder
	GetConfig() api.Config
}
