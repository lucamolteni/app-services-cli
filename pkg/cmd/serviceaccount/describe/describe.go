package describe

import (
	"context"

	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/config"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/iostreams"
	"github.com/redhat-developer/app-services-cli/pkg/core/localize"
	"github.com/redhat-developer/app-services-cli/pkg/core/logging"
	svcacctmgmterrors "github.com/redhat-developer/app-services-cli/pkg/shared/svcacctmgmtutil"

	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"

	"github.com/spf13/cobra"
)

type options struct {
	id           string
	outputFormat string
	enableAuthV2 bool

	IO         *iostreams.IOStreams
	Config     config.IConfig
	Connection factory.ConnectionFunc
	localizer  localize.Localizer
	Context    context.Context
	Logger     logging.Logger
}

func NewDescribeCommand(f *factory.Factory) *cobra.Command {
	opts := &options{
		Config:     f.Config,
		Connection: f.Connection,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
		Context:    f.Context,
		Logger:     f.Logger,
	}

	cmd := &cobra.Command{
		Use:     "describe",
		Short:   opts.localizer.MustLocalize("serviceAccount.describe.cmd.shortDescription"),
		Long:    opts.localizer.MustLocalize("serviceAccount.describe.cmd.longDescription"),
		Example: opts.localizer.MustLocalize("serviceAccount.describe.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			validOutputFormats := flagutil.ValidOutputFormats
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, validOutputFormats...) {
				return flagutil.InvalidValueError("output", opts.outputFormat, validOutputFormats...)
			}

			return runDescribe(opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", opts.localizer.MustLocalize("serviceAccount.describe.flag.id.description"))
	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "json", opts.localizer.MustLocalize("serviceAccount.common.flag.output.description"))
	cmd.Flags().BoolVar(&opts.enableAuthV2, "enable-auth-v2", false, opts.localizer.MustLocalize("serviceAccount.common.flag.enableAuthV2"))

	_ = cmd.MarkFlagRequired("id")

	_ = cmd.Flags().MarkDeprecated("enable-auth-v2", opts.localizer.MustLocalize("serviceAccount.common.flag.deprecated.enableAuthV2"))

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

func runDescribe(opts *options) error {
	conn, err := opts.Connection()
	if err != nil {
		return err
	}

	api := conn.KiotaAPI()

	res, err := api.ServiceAccountMgmt().V1ById(opts.id).Get(opts.Context, nil)

	if apiErr := svcacctmgmterrors.GetAPIError(err); apiErr != nil {
		switch apiErr.GetError() {
		case "service_account_not_found":
			return opts.localizer.MustLocalizeError("serviceAccount.common.error.notFoundError", localize.NewEntry("ID", opts.id))
		default:
			return err
		}
	}

	return dump.Formatted(opts.IO.Out, opts.outputFormat, res)
}
