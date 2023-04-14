package flagutil

import (
	kiotapi "github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt/api"
	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/localize"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type FlagSet struct {
	flags     *pflag.FlagSet
	cmd       *cobra.Command
	localizer localize.Localizer
	*flagutil.FlagSet
}

// NewFlagSet returns a new flag set for creating common Kafka-command flags
func NewFlagSet(cmd *cobra.Command, localizer localize.Localizer) *FlagSet {
	return &FlagSet{
		cmd:       cmd,
		flags:     cmd.Flags(),
		localizer: localizer,
		FlagSet:   flagutil.NewFlagSet(cmd, localizer),
	}
}

// AddInstanceID adds a flag for setting the Kafka instance ID
func (fs *FlagSet) AddInstanceID(instanceID *string) {
	flagName := "instance-id"

	fs.flags.StringVar(
		instanceID,
		flagName,
		"",
		flagutil.FlagDescription(fs.localizer, "kafka.common.flag.instanceID.description"),
	)
}

// RegisterNameFlagCompletionFunc adds dynamic completion for the --name flag
func RegisterNameFlagCompletionFunc(cmd *cobra.Command, f *factory.Factory) error {
	return cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var validNames []string
		directive := cobra.ShellCompDirectiveNoSpace

		conn, err := f.Connection()
		if err != nil {
			return validNames, directive
		}

		var searchQ *string
		if toComplete != "" {
			queryString := "name like " + toComplete + "%"
			searchQ = &queryString

		}

		kafkas, err := conn.KiotaAPI().KafkaMgmt().V1().Kafkas().Get(f.Context, &kiotapi.Kafkas_mgmtV1KafkasRequestBuilderGetRequestConfiguration{
			QueryParameters: &kiotapi.Kafkas_mgmtV1KafkasRequestBuilderGetQueryParameters{
				Search: searchQ,
			},
		})

		if err != nil {
			return validNames, directive
		}

		items := kafkas.GetItems()
		for index := range items {
			validNames = append(validNames, *items[index].GetName())
		}

		return validNames, directive
	})
}
