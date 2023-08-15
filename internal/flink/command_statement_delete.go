package flink

import (
	"fmt"

	"github.com/spf13/cobra"

	pcmd "github.com/confluentinc/cli/v3/pkg/cmd"
	"github.com/confluentinc/cli/v3/pkg/errors"
	"github.com/confluentinc/cli/v3/pkg/form"
	"github.com/confluentinc/cli/v3/pkg/output"
	"github.com/confluentinc/cli/v3/pkg/resource"
)

func (c *command) newStatementDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete <name>",
		Short:             "Delete a Flink SQL statement.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: pcmd.NewValidArgsFunction(c.validStatementArgs),
		RunE:              c.statementDelete,
	}

	pcmd.AddCloudFlag(cmd)
	c.addRegionFlag(cmd)
	pcmd.AddForceFlag(cmd)
	pcmd.AddEnvironmentFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddContextFlag(cmd, c.CLICommand)

	return cmd
}

func (c *command) statementDelete(cmd *cobra.Command, args []string) error {
	environmentId, err := c.Context.EnvironmentId()
	if err != nil {
		return err
	}

	client, err := c.GetFlinkGatewayClient()
	if err != nil {
		return err
	}

	if _, err := client.GetStatement(environmentId, args[0], c.Context.LastOrgId); err != nil {
		return err
	}

	promptMsg := fmt.Sprintf(errors.DeleteResourceConfirmYesNoMsg, resource.FlinkStatement, args[0])
	if ok, err := form.ConfirmDeletion(cmd, promptMsg, ""); err != nil || !ok {
		return err
	}

	if err := client.DeleteStatement(environmentId, args[0], c.Context.LastOrgId); err != nil {
		return err
	}

	output.Printf(errors.DeletedResourceMsg, resource.FlinkStatement, args[0])
	return nil
}