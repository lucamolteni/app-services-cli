package openshiftclustercmdutil

import (
	"fmt"
	"github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	kmodels "github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt/models"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/kafka/flagutil"
	"strconv"

	"github.com/redhat-developer/app-services-cli/pkg/core/errors"
	"github.com/redhat-developer/app-services-cli/pkg/core/localize"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
)

// Validator is a type for validating Kafka configuration values
type Validator struct {
	Localizer  localize.Localizer
	Connection factory.ConnectionFunc
}

func ValidateMachinePoolCount(count int) bool {
	// check if the count is a multiple of 3 and greater than or equal to 3
	if count%3 == 0 && count >= 3 {
		return true
	}
	return false
}

func (v *Validator) ValidatorForMachinePoolNodes(val interface{}) error {
	value := fmt.Sprintf("%v", val)
	if val == "" {
		return errors.NewCastError(val, "empty string")
	}
	value1, err := strconv.Atoi(value)
	if err != nil {
		return errors.NewCastError(val, "integer")
	}
	if !ValidateMachinePoolCount(value1) {
		return fmt.Errorf("invalid input, machine pool node count must be greater than or equal to 3 and it " +
			"must be a is a multiple of 3")
	}
	return nil
}

func CreatePromptOptionsFromClusters(clusterList *kmodels.EnterpriseClusterList, clusterMap *map[string]v1.Cluster) *[]string {
	promptOptions := []string{}
	validatedClusters := ValidateClusters(clusterList)
	for _, cluster := range validatedClusters.GetItems() {
		ocmCluster := (*clusterMap)[*cluster.GetId()]
		display := ocmCluster.Name() + " (" + *cluster.GetId() + ")"
		promptOptions = append(promptOptions, display)
	}
	return &promptOptions
}

func ValidateClusters(clusterList *kmodels.EnterpriseClusterList) *kmodels.EnterpriseClusterList {
	// if cluster is in a ready state add it to the list of clusters
	//items := make([]kmodels.EnterpriseClusterListItemable, 0, len(clusterList.GetItems()))
	//for _, cluster := range clusterList.GetItems() {
	//	if *cluster.GetStatus() == "ready" {
	//		items = append(items, cluster)
	//	}
	//}
	//
	////newClusterList := kafkamgmtclient.NewEnterpriseClusterList(clusterList.Kind, clusterList.Page, int32(len(items)), clusterList.Total, items)
	//newClusterList := kmodels.NewEnterpriseClusterList()
	//newClusterList.SetItems(items)
	//return newClusterList

	// TODO EnterpriseClusterListItemable doesn't exist

	return nil
}

func HideClusterMgmtFlags(flags *flagutil.FlagSet) {
	_ = flags.MarkHidden("cluster-mgmt-api-url")
	_ = flags.MarkHidden("access-token")
}
