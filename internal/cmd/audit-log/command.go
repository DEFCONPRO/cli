package auditlog

import (
	"encoding/json"
	"fmt"
	"net/http"

	mds "github.com/confluentinc/mds-sdk-go/mdsv1"
	"github.com/spf13/cobra"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/errors"
)

type command struct {
	*pcmd.CLICommand
	prerunner pcmd.PreRunner
}

// New returns the default command object for interacting with audit logs.
func New(prerunner pcmd.PreRunner) *cobra.Command {
	cmd := &cobra.Command{
		Use:         "audit-log",
		Short:       "Manage audit log configuration.",
		Long:        "Manage which auditable events are logged, and where the event logs are sent.",
		Annotations: map[string]string{pcmd.RunRequirement: pcmd.RequireCloudLoginOrOnPremLogin},
	}

	c := &command{
		CLICommand: pcmd.NewAnonymousCLICommand(cmd, prerunner),
		prerunner:  prerunner,
	}
	c.init()

	return c.Command
}

func (c *command) init() {
	c.AddCommand(NewDescribeCommand(c.prerunner))
	c.AddCommand(NewMigrateCommand(c.prerunner))
	c.AddCommand(NewConfigCommand(c.prerunner))
	c.AddCommand(NewRouteCommand(c.prerunner))
}

type errorMessage struct {
	ErrorCode uint32 `json:"error_code" yaml:"error_code"`
	Message   string `json:"message" yaml:"message"`
}

func HandleMdsAuditLogApiError(cmd *cobra.Command, err error, response *http.Response) error {
	if response != nil {
		switch status := response.StatusCode; status {
		case http.StatusNotFound:
			cmd.SilenceUsage = true
			return errors.NewWrapErrorWithSuggestions(err, errors.UnableToAccessEndpointErrorMsg, errors.UnableToAccessEndpointSuggestions)
		case http.StatusForbidden:
			switch e := err.(type) {
			case mds.GenericOpenAPIError:
				cmd.SilenceUsage = true
				em := errorMessage{}
				if err = json.Unmarshal(e.Body(), &em); err != nil {
					// It wasn't what we expected. Use the regular error handler.
					return errors.HandleCommon(err, cmd)
				}
				return fmt.Errorf("%s\n%s", e.Error(), em.Message)
			}
		}
	}
	return err
}