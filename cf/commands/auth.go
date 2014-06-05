package commands

import (
	"errors"
	"github.com/starkandwayne/cf-cli/cf"
	"github.com/starkandwayne/cf-cli/cf/api"
	"github.com/starkandwayne/cf-cli/cf/command_metadata"
	"github.com/starkandwayne/cf-cli/cf/configuration"
	"github.com/starkandwayne/cf-cli/cf/requirements"
	"github.com/starkandwayne/cf-cli/cf/terminal"
	"github.com/codegangsta/cli"
)

type Authenticate struct {
	ui            terminal.UI
	config        configuration.ReadWriter
	authenticator api.AuthenticationRepository
}

func NewAuthenticate(ui terminal.UI, config configuration.ReadWriter, authenticator api.AuthenticationRepository) (cmd Authenticate) {
	cmd.ui = ui
	cmd.config = config
	cmd.authenticator = authenticator
	return
}

func (cmd Authenticate) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "auth",
		Description: "Authenticate user non-interactively",
		Usage: "CF_NAME auth USERNAME PASSWORD\n\n" +
			terminal.WarningColor("WARNING:\n   Providing your password as a command line option is highly discouraged\n   Your password may be visible to others and may be recorded in your shell history\n\n") +
			"EXAMPLE:\n" +
			"   CF_NAME auth name@example.com \"my password\" (use quotes for passwords with a space)\n" +
			"   CF_NAME auth name@example.com \"\\\"password\\\"\" (escape quotes if used in password)",
	}
}

func (cmd Authenticate) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 2 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c)
		return
	}

	reqs = append(reqs, requirementsFactory.NewApiEndpointRequirement())
	return
}

func (cmd Authenticate) Run(c *cli.Context) {
	cmd.config.ClearSession()
	cmd.authenticator.GetLoginPromptsAndSaveUAAServerURL()

	cmd.ui.Say("API endpoint: %s", terminal.EntityNameColor(cmd.config.ApiEndpoint()))
	cmd.ui.Say("Authenticating...")

	apiErr := cmd.authenticator.Authenticate(map[string]string{
		"username": c.Args()[0],
		"password": c.Args()[1],
	})
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
	cmd.ui.Say("Use '%s' to view or set your target org and space", terminal.CommandColor(cf.Name()+" target"))
	return
}
