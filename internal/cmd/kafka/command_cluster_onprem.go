package kafka

import (
	"context"
	"github.com/spf13/cobra"

	"github.com/antihax/optional"
	print "github.com/confluentinc/cli/internal/pkg/cluster"
	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/output"
	mds "github.com/confluentinc/mds-sdk-go/mdsv1"
)

var kafkaClusterTypeName = "kafka-cluster"

type clusterCommandOnPrem struct {
	*pcmd.AuthenticatedCLICommand
	prerunner pcmd.PreRunner
}

// NewClusterCommand returns the Cobra command for Kafka cluster.
func NewClusterCommandOnPrem(prerunner pcmd.PreRunner) *cobra.Command {
	cliCmd := pcmd.NewAuthenticatedWithMDSCLICommand(
		&cobra.Command{
			Use:   "cluster",
			Short: "Manage Kafka clusters.",
		},
		prerunner)
	cmd := &clusterCommandOnPrem{
		AuthenticatedCLICommand: cliCmd,
		prerunner:               prerunner,
	}
	cmd.init()
	return cmd.Command
}

func (c *clusterCommandOnPrem) init() {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered Kafka clusters.",
		Long:  "List Kafka clusters that are registered with the MDS cluster registry.",
		RunE:  c.list,
		Args:  cobra.NoArgs,
	}
	listCmd.Flags().StringP(output.FlagName, output.ShortHandFlag, output.DefaultValue, output.Usage)
	listCmd.Flags().SortFlags = false
	c.AddCommand(listCmd)
}

func (c *clusterCommandOnPrem) createContext() context.Context {
	return context.WithValue(context.Background(), mds.ContextAccessToken, c.State.AuthToken)
}

func (c *clusterCommandOnPrem) list(cmd *cobra.Command, args []string) error {
	clustertype := &mds.ClusterRegistryListOpts{
		ClusterType: optional.NewString(kafkaClusterTypeName),
	}
	clusterInfos, response, err := c.MDSClient.ClusterRegistryApi.ClusterRegistryList(c.createContext(), clustertype)
	if err != nil {
		return print.HandleClusterError(cmd, err, response)
	}
	format, err := cmd.Flags().GetString(output.FlagName)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	return print.PrintCluster(cmd, clusterInfos, format)
}