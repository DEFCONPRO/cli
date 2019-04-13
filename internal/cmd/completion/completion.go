package completion

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/spf13/cobra"
)

const longDescriptionTemplate = `Output shell completion code for the specified shell (bash only).
The shell code must be evaluated to provide interactive completion of {{.CLIName}} commands.

Install bash completions on macOS:

  # Enable bash completions using homebrew
  brew install bash-completion
  # Set the {{.CLIName}} completion code for bash to a file that's sourced on login
  {{.CLIName}} completion bash > $(brew --prefix)/etc/bash_completion.d/{{.CLIName}}

Install bash completions on Linux:

  # Set the {{.CLIName}} completion code for bash to a file that's sourced on login
  {{.CLIName}} completion bash > /etc/bash_completion.d/{{.CLIName}}

  # Load the {{.CLIName}} completion code for bash into the current shell
  source /etc/bash_completion.d/{{.CLIName}}
`

type completionCommand struct {
	*cobra.Command
	rootCmd *cobra.Command
}

// NewCompletionCmd returns the Cobra command for shell completion.
func NewCompletionCmd(rootCmd *cobra.Command, cliName string) (*cobra.Command, error) {
	cmd := &completionCommand{
		rootCmd: rootCmd,
	}
	err := cmd.init(cliName)
	return cmd.Command, err
}

func (c *completionCommand) init(cliName string) error {
	longDescription, err := getLongDescription(cliName)
	if err != nil {
		return err
	}

	c.Command = &cobra.Command{
		Use:   "completion SHELL",
		Short: "Output shell completion code",
		Long:  longDescription,
		RunE:  c.completion,
		Args:  cobra.ExactArgs(1),
	}
	return nil
}

func (c *completionCommand) completion(cmd *cobra.Command, args []string) error {
	var err error
	if args[0] == "bash" {
		err = c.rootCmd.GenBashCompletion(cmd.OutOrStdout())
	} else {
		err = fmt.Errorf(`unsupported shell type "%s"`, args[0])
	}
	return err
}

func getLongDescription(cliName string) (string, error) {
	t := template.Must(template.New("longDescription").Parse(longDescriptionTemplate))
	buf := new(bytes.Buffer)
	data := map[string]interface{}{"CLIName": cliName}
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
