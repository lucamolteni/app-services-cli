package promote

import (
	"github.com/redhat-developer/app-services-cli/pkg/cmd/kafka/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/localize"
	"github.com/redhat-developer/app-services-cli/pkg/shared/contextutil"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/redhat-developer/app-services-cli/pkg/shared/kafkautil"
	"github.com/spf13/cobra"

	"github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt/models"
	kafkamgmtv1errors "github.com/redhat-developer/app-services-cli/pkg/shared/kafkautil/errors"
)

type options struct {
	id                  string
	name                string
	marketplaceAcctId   string
	marketplace         string
	desiredBillingModel string

	f *factory.Factory
}

func NewPromoteCommand(f *factory.Factory) *cobra.Command {

	opts := &options{
		f: f,
	}

	cmd := &cobra.Command{
		Use:     "promote",
		Short:   opts.f.Localizer.MustLocalize("kafka.promote.cmd.shortDescription"),
		Long:    opts.f.Localizer.MustLocalize("kafka.promote.cmd.longDescription"),
		Example: opts.f.Localizer.MustLocalize("kafka.promote.cmd.example"),
		Args:    cobra.NoArgs,
		Hidden:  true,
		RunE: func(cmd *cobra.Command, args []string) error {

			if opts.name != "" && opts.id != "" {
				return opts.f.Localizer.MustLocalizeError("service.error.idAndNameCannotBeUsed")
			}

			if opts.id != "" || opts.name != "" {
				return runPromote(opts)
			}

			kafkaInstance, err := contextutil.GetCurrentKafkaInstance(f)
			if err != nil {
				return err
			}

			opts.id = kafkaInstance.GetId()

			return runPromote(opts)
		},
	}

	flags := flagutil.NewFlagSet(cmd, opts.f.Localizer)

	flags.StringVar(&opts.id, "id", "", opts.f.Localizer.MustLocalize("kafka.promote.flag.id"))
	flags.StringVar(&opts.name, "name", "", opts.f.Localizer.MustLocalize("kafka.promote.flag.name"))

	flags.StringVar(&opts.marketplaceAcctId, "marketplace-account-id", "", f.Localizer.MustLocalize("kafka.common.flag.marketplaceId.description"))
	flags.StringVar(&opts.marketplace, "marketplace", "", f.Localizer.MustLocalize("kafka.common.flag.marketplaceType.description"))
	flags.StringVar(&opts.desiredBillingModel, "billing-model", "", f.Localizer.MustLocalize("kafka.common.flag.billingModel.description"))

	_ = cmd.MarkFlagRequired("billing-model")

	return cmd

}

func runPromote(opts *options) error {

	conn, err := opts.f.Connection()
	if err != nil {
		return err
	}

	api := conn.KiotaAPI()

	if opts.name != "" {
		response, _, newErr := kafkautil.GetKafkaByNameK(opts.f.Context, api.KafkaMgmt(), opts.name)
		if newErr != nil {
			return newErr
		}

		opts.id = *(*response).GetId()
	}

	a := api.KafkaMgmt().V1().KafkasById(opts.id).Promote()

	var promoteOptions models.KafkaPromoteRequest

	promoteOptions.SetDesiredKafkaBillingModel(&opts.desiredBillingModel)

	if opts.marketplace != "" {
		promoteOptions.SetDesiredMarketplace(&opts.marketplace)
	}

	if opts.marketplaceAcctId != "" {
		promoteOptions.SetDesiredKafkaBillingModel(&opts.marketplaceAcctId)
	}

	err = a.Post(opts.f.Context, &promoteOptions, nil)

	if apiErr := kafkautil.GetAPIError(err); apiErr != nil {
		switch apiErr.GetCode() {
		case kafkamgmtv1errors.ERROR_120:
			// For standard instances
			return opts.f.Localizer.MustLocalizeError("kafka.create.error.quota.exceeded")
		case kafkamgmtv1errors.ERROR_24:
			// For dev instances
			return opts.f.Localizer.MustLocalizeError("kafka.create.error.instance.limit")
		case kafkamgmtv1errors.ERROR_9:
			return opts.f.Localizer.MustLocalizeError("kafka.create.error.standard.promote")
		case kafkamgmtv1errors.ERROR_43:
			return opts.f.Localizer.MustLocalizeError("kafka.create.error.billing.invalid", localize.NewEntry("Billing", opts.marketplaceAcctId))
		}
	}

	if err != nil {
		return err
	}

	opts.f.Logger.Info(opts.f.Localizer.MustLocalize("kafka.promote.info.successAsync"))

	return nil
}
