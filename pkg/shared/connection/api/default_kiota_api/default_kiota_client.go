package default_kiota_api

import (
	"context"
	"fmt"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	khttp "github.com/microsoft/kiota-http-go"
	"net/url"

	"github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt"
	kafkamgmtapi "github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt/api"
	"github.com/redhat-developer/app-services-cli/pkg/shared/connection/api"

	svcacctmgmtclient "github.com/redhat-developer/app-services-sdk-core/app-services-sdk-go/serviceaccountmgmt/apiv1/client"
)

// defaultAPI is a type which defines a number of API creator functions
type defaultAPI struct {
	api.Config
}

type RedHatAccessTokenProvider struct {
	accessToken string
}

func (r RedHatAccessTokenProvider) GetAuthorizationToken(context context.Context, url *url.URL, additionalAuthenticationContext map[string]interface{}) (string, error) {
	return r.accessToken, nil
}

func (r RedHatAccessTokenProvider) GetAllowedHostsValidator() *authentication.AllowedHostsValidator {
	return nil
}

// New creates a new default API client wrapper
func New(cfg *api.Config) *defaultAPI {
	return &defaultAPI{
		Config: *cfg,
	}
}

func (a *defaultAPI) GetConfig() api.Config {
	return a.Config
}

// KafkaMgmt returns a new Kafka Management API client instance
func (a *defaultAPI) KafkaMgmt() *kafkamgmtapi.ApiRequestBuilder {

	tokenProvider := RedHatAccessTokenProvider{accessToken: a.GetConfig().AccessToken}

	provider := authentication.NewBaseBearerTokenAuthenticationProvider(tokenProvider)

	options := khttp.ObservabilityOptions{}
	options.SetIncludeEUIIAttributes(true)

	adapter, err := khttp.NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClientAndObservabilityOptions(provider, nil, nil, nil, options)

	if err != nil {
		fmt.Printf("Error creating request adapter: %v\n", err)
	}

	//BaseURL:    a.ApiURL.String(),
	//Debug:      a.Logger.DebugEnabled(),
	//HTTPClient: tc,
	//UserAgent:  a.UserAgent,

	kiotaClient := kafkamgmt.NewApiClient(adapter)

	kiotaAPI := kiotaClient.Api()

	return kiotaAPI
}

// ServiceAccountMgmt return a new Service Account Management API client instance
func (a *defaultAPI) ServiceAccountMgmt() svcacctmgmtclient.ServiceAccountsApi {

	// TODO this is needed for  service account create

	return nil
	//tc := a.CreateOAuthTransport(a.AccessToken)
	//client := svcacctmgmt.NewAPIClient(&svcacctmgmt.Config{
	//	BaseURL:    a.AuthURL.String(),
	//	Debug:      a.Logger.DebugEnabled(),
	//	HTTPClient: tc,
	//	UserAgent:  a.UserAgent,
	//})
	//
	//return client.ServiceAccountsApi
}
