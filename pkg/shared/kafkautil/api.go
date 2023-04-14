package kafkautil

import (
	"context"
	"fmt"
	"github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt/api"
	"github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt/models"
	"net/http"

	kafkapi "github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt/api"
	kafkamgmtv1errors "github.com/redhat-developer/app-services-cli/pkg/shared/kafkautil/errors"
	kafkamgmtclient "github.com/redhat-developer/app-services-sdk-core/app-services-sdk-go/kafkamgmt/apiv1/client"
)

func GetKafkaByID(ctx context.Context, api kafkamgmtclient.DefaultApi, id string) (*kafkamgmtclient.KafkaRequest, *http.Response, error) {
	r := api.GetKafkaById(ctx, id)

	kafkaReq, httpResponse, err := r.Execute()
	if IsAPIError(err, kafkamgmtv1errors.ERROR_7) {
		return nil, httpResponse, NotFoundByIDError(id)
	}

	return &kafkaReq, httpResponse, err
}

func GetKafkaByIDK(ctx context.Context, api *api.Kafkas_mgmtRequestBuilder, id string) (*models.KafkaRequestable, *http.Response, error) {

	kafkaReq, err := api.V1().KafkasById(id).Get(ctx, nil)
	if IsAPIError(err, kafkamgmtv1errors.ERROR_7) {
		return nil, nil, NotFoundByIDError(id)
	}

	return &kafkaReq, nil, err
}

func GetKafkaByName(ctx context.Context, api kafkamgmtclient.DefaultApi, name string) (*kafkamgmtclient.KafkaRequest, *http.Response, error) {
	r := api.GetKafkas(ctx)
	r = r.Search(fmt.Sprintf("name = %v", name))
	kafkaList, httpResponse, err := r.Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	if kafkaList.GetTotal() == 0 {
		return nil, nil, NotFoundByNameError(name)
	}

	items := kafkaList.GetItems()
	kafkaReq := items[0]

	return &kafkaReq, httpResponse, err
}

func GetKafkaByNameK(ctx context.Context, api *kafkapi.Kafkas_mgmtRequestBuilder, name string) (*models.KafkaRequestable, *http.Response, error) {
	queryString := fmt.Sprintf("name = %v", name)
	kafkaList, err := api.V1().Kafkas().Get(ctx, &kafkapi.Kafkas_mgmtV1KafkasRequestBuilderGetRequestConfiguration{
		QueryParameters: &kafkapi.Kafkas_mgmtV1KafkasRequestBuilderGetQueryParameters{
			Search: &queryString,
		},
	})

	if err != nil {
		return nil, nil, err
	}

	if *kafkaList.GetTotal() == 0 {
		return nil, nil, NotFoundByNameError(name)
	}

	items := kafkaList.GetItems()
	kafkaReq := items[0]

	return &kafkaReq, nil, err
}
