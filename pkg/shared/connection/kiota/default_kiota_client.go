package kiota

import (
	"context"
	"fmt"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	khttp "github.com/microsoft/kiota-http-go"
	nethttp "net/http"
	"net/url"

	"github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt"
	kafkamgmtapi "github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt/api"
	"github.com/redhat-developer/app-services-cli/pkg/shared/connection/api"

	"github.com/redhat-developer/app-services-cli/pkg/apisdk/svcacctmgmt"
	svcacctmgmtapis "github.com/redhat-developer/app-services-cli/pkg/apisdk/svcacctmgmt/apis"
	"github.com/redhat-developer/app-services-cli/pkg/shared/kiotalog"
)

// defaultKiotaAPIClient is a type which defines a number of API creator functions
type defaultKiotaAPIClient struct {
	config api.Config
}

func (a *defaultKiotaAPIClient) GetConfig() api.Config {
	return a.config
}

type RedHatAccessTokenProvider struct {
	accessToken string
}

func (r RedHatAccessTokenProvider) GetAuthorizationToken(context.Context, *url.URL, map[string]interface{}) (string, error) {
	return r.accessToken, nil
}

func (r RedHatAccessTokenProvider) GetAllowedHostsValidator() *authentication.AllowedHostsValidator {
	return nil
}

// New creates a new default API client wrapper
func New(cfg *api.Config) *defaultKiotaAPIClient {
	return &defaultKiotaAPIClient{
		config: *cfg,
	}
}

func (client *defaultKiotaAPIClient) adapter() *khttp.NetHttpRequestAdapter {

	tokenProvider := RedHatAccessTokenProvider{accessToken: client.config.AccessToken}
	authenticationProvider := authentication.NewBaseBearerTokenAuthenticationProvider(tokenProvider)

	httpClient := client.createHttpClient()

	adapter, err := khttp.NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClientAndObservabilityOptions(
		authenticationProvider,
		nil,
		nil,
		httpClient,
		khttp.ObservabilityOptions{})

	if err != nil {
		fmt.Printf("Error creating request adapter: %v\n", err)
	}
	return adapter
}

func (client *defaultKiotaAPIClient) createHttpClient() *nethttp.Client {
	var httpClient *nethttp.Client
	if client.config.Logger.DebugEnabled() {
		middlewares := khttp.GetDefaultMiddlewares()
		logger := kiotalog.NewLoggingHandler(client.config.Logger)
		withLogger := append(middlewares, logger)
		httpClient = khttp.GetDefaultClient(withLogger...)
	} else {
		httpClient = khttp.GetDefaultClient()
	}
	return httpClient
}

// KafkaMgmt returns a new Kafka Management API client instance
func (client *defaultKiotaAPIClient) KafkaMgmt() *kafkamgmtapi.Kafkas_mgmtRequestBuilder {
	return kafkamgmt.NewApiClient(client.adapter()).Api().Kafkas_mgmt()
}

// ServiceAccountMgmt return a new Service Account Management API client instance
func (client *defaultKiotaAPIClient) ServiceAccountMgmt() *svcacctmgmtapis.Service_accountsRequestBuilder {
	return svcacctmgmt.NewApiClient(client.adapter()).Apis().Service_accounts()
}
