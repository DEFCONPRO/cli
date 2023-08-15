package kafka

import (
	"github.com/spf13/cobra"

	pcmd "github.com/confluentinc/cli/v3/pkg/cmd"
	"github.com/confluentinc/cli/v3/pkg/examples"
	"github.com/confluentinc/cli/v3/pkg/output"
)

func (c *consumerGroupCommand) newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Kafka consumer groups.",
		Args:  cobra.NoArgs,
		RunE:  c.list,
		Example: examples.BuildExampleString(
			examples.Example{
				Text: "List all consumer groups.",
				Code: "confluent kafka consumer-group list",
			},
		),
		Hidden: true,
	}

	pcmd.AddClusterFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddContextFlag(cmd, c.CLICommand)
	pcmd.AddEnvironmentFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddOutputFlag(cmd)

	return cmd
}

func (c *consumerGroupCommand) list(cmd *cobra.Command, _ []string) error {
	kafkaREST, err := c.GetKafkaREST()
	if err != nil {
		return err
	}

	groupCmdResp, err := kafkaREST.CloudClient.ListKafkaConsumerGroups()
	if err != nil {
		return err
	}

	list := output.NewList(cmd)
	for _, group := range groupCmdResp.Data {
		list.Add(&consumerGroupOut{
			ClusterId:         group.GetClusterId(),
			ConsumerGroupId:   group.GetConsumerGroupId(),
			Coordinator:       getStringBroker(group.GetCoordinator()),
			IsSimple:          group.GetIsSimple(),
			PartitionAssignor: group.GetPartitionAssignor(),
			State:             group.GetState(),
		})
	}
	return list.Print()
}