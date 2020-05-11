package connector_catalog

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	connectv1 "github.com/confluentinc/ccloudapis/connect/v1"
	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	v3 "github.com/confluentinc/cli/internal/pkg/config/v3"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/output"
)

type command struct {
	*pcmd.AuthenticatedCLICommand
}

type catalogDisplay struct {
	PluginName string
	Type       string
}

var (
	catalogFields          = []string{"PluginName", "Type"}
	catalogStructureLabels = []string{"plugin_name", "type"}
)

// New returns the default command object for interacting with Connect.
func New(prerunner pcmd.PreRunner, config *v3.Config) *cobra.Command {
	cmd := &command{
		AuthenticatedCLICommand: pcmd.NewAuthenticatedCLICommand(&cobra.Command{
			Use:   "connector-catalog",
			Short: "Catalog of connectors and their configurations.",
		}, config, prerunner),
	}
	cmd.init()
	return cmd.Command
}

func (c *command) init() {
	cmd := &cobra.Command{
		Use:   "describe <connector-type>",
		Short: "Describe a connector plugin type.",
		Example: FormatDescription(`
Describe required connector configuration parameters for a specific connector plugin.
With the --sample-file flag, create a sample connector configuration file.
::

        {{.CLIName}} connector-catalog describe <PluginName>
        {{.CLIName}} connector-catalog describe <PluginName> --sample-file <filename>`, c.Config.CLIName),
		RunE: c.describe,
		Args: cobra.ExactArgs(1),
	}
	cmd.Flags().String("cluster", "", "Kafka cluster ID.")
	cmd.Flags().StringP(output.FlagName, output.ShortHandFlag, output.DefaultValue, output.Usage)
	cmd.Flags().SortFlags = false
	c.AddCommand(cmd)

	cmd = &cobra.Command{
		Use:   "list",
		Short: "List connector plugin types.",
		Example: FormatDescription(`
List connectors in the current or specified Kafka cluster context.

::

        {{.CLIName}} connector-catalog list`, c.Config.CLIName),
		RunE: c.list,
		Args: cobra.NoArgs,
	}
	cmd.Flags().String("cluster", "", "Kafka cluster ID.")
	cmd.Flags().StringP(output.FlagName, output.ShortHandFlag, output.DefaultValue, output.Usage)
	cmd.Flags().SortFlags = false
	c.AddCommand(cmd)
}

func (c *command) list(cmd *cobra.Command, args []string) error {
	kafkaCluster, err := c.Context.GetKafkaClusterForCommand(cmd)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	connectorInfo, err := c.Client.Connect.GetPlugins(context.Background(), &connectv1.Connector{AccountId: c.EnvironmentId(), KafkaClusterId: kafkaCluster.ID}, "")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	outputWriter, err := output.NewListOutputWriter(cmd, catalogFields, catalogFields, catalogStructureLabels)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	for _, conn := range connectorInfo {
		connector := &catalogDisplay{
			PluginName: conn.Class,
			Type:       conn.Type,
		}
		outputWriter.AddElement(connector)
	}
	return outputWriter.Out()
}

func (c *command) describe(cmd *cobra.Command, args []string) error {
	kafkaCluster, err := c.Context.GetKafkaClusterForCommand(cmd)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	if len(args) == 0 {
		return errors.HandleCommon(errors.ErrNoPluginName, cmd)
	}
	config := map[string]string{"connector.class": args[0]}

	reply, err := c.Client.Connect.Validate(context.Background(),
		&connectv1.ConnectorConfig{
			UserConfigs:    config,
			AccountId:      c.EnvironmentId(),
			KafkaClusterId: kafkaCluster.ID,
			Plugin:         args[0]})
	if reply != nil && err != nil {
		outputFormat, flagErr := cmd.Flags().GetString(output.FlagName)
		if flagErr != nil {
			return errors.HandleCommon(flagErr, cmd)
		}
		if outputFormat == output.Human.String() {
			pcmd.Println(cmd, "Following are the required configs: \nconnector.class: "+args[0]+"\n"+err.Error())
		} else {

			for _, c := range reply.Configs {
				if len(c.Value.Errors) > 0 {
					config[c.Value.Name] = fmt.Sprintf("%s ", c.Value.Errors[0])
				}
			}
			return output.StructuredOutput(outputFormat, &config)
		}
		return nil
	}
	return errors.HandleCommon(errors.ErrInvalidCloud, cmd)
}

func FormatDescription(description string, cliName string) string {
	return strings.ReplaceAll(description, "{{.CLIName}}", cliName)
}
