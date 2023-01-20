package schemaregistry

import (
	"context"

	srsdk "github.com/confluentinc/schema-registry-sdk-go"
	"github.com/spf13/cobra"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/output"
)

type listOut struct {
	Exporter string `human:"Exporter" serialized:"exporter"`
}

func (c *exporterCommand) newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all schema exporters.",
		Args:  cobra.NoArgs,
		RunE:  c.list,
	}

	pcmd.AddApiKeyFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddApiSecretFlag(cmd)
	pcmd.AddContextFlag(cmd, c.CLICommand)
	pcmd.AddEnvironmentFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddOutputFlag(cmd)

	return cmd
}

func (c *exporterCommand) list(cmd *cobra.Command, _ []string) error {
	srClient, ctx, err := getApiClient(cmd, c.srClient, c.Config, c.Version)
	if err != nil {
		return err
	}

	return listExporters(cmd, srClient, ctx)
}

func listExporters(cmd *cobra.Command, srClient *srsdk.APIClient, ctx context.Context) error {
	exporters, _, err := srClient.DefaultApi.GetExporters(ctx)
	if err != nil {
		return err
	}

	list := output.NewList(cmd)
	for _, exporter := range exporters {
		list.Add(&listOut{Exporter: exporter})
	}
	return list.Print()
}
