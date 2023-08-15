package organization

import (
	"github.com/spf13/cobra"

	orgv2 "github.com/confluentinc/ccloud-sdk-go-v2/org/v2"

	pcmd "github.com/confluentinc/cli/v3/pkg/cmd"
	"github.com/confluentinc/cli/v3/pkg/errors"
	"github.com/confluentinc/cli/v3/pkg/output"
	"github.com/confluentinc/cli/v3/pkg/resource"
)

func (c *command) newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the current Confluent Cloud organization.",
		Args:  cobra.NoArgs,
		RunE:  c.update,
	}

	cmd.Flags().String("name", "", "Name of the Confluent Cloud organization.")
	pcmd.AddOutputFlag(cmd)

	cobra.CheckErr(cmd.MarkFlagRequired("name"))

	return cmd
}

func (c *command) update(cmd *cobra.Command, args []string) error {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	organization := orgv2.OrgV2Organization{DisplayName: orgv2.PtrString(name)}
	organization, httpResp, err := c.V2Client.UpdateOrgOrganization(c.Context.GetCurrentOrganization(), organization)
	if err != nil {
		return errors.CatchCCloudV2ResourceNotFoundError(err, resource.Organization, httpResp)
	}

	table := output.NewTable(cmd)
	table.Add(&out{
		IsCurrent: organization.GetId() == c.Context.GetCurrentOrganization(),
		Id:        organization.GetId(),
		Name:      organization.GetDisplayName(),
	})
	return table.Print()
}