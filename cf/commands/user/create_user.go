package user

import (
	"github.com/starkandwayne/cf-cli/cf"
	"github.com/starkandwayne/cf-cli/cf/api"
	"github.com/starkandwayne/cf-cli/cf/command_metadata"
	"github.com/starkandwayne/cf-cli/cf/configuration"
	"github.com/starkandwayne/cf-cli/cf/errors"
	"github.com/starkandwayne/cf-cli/cf/requirements"
	"github.com/starkandwayne/cf-cli/cf/terminal"
	"github.com/codegangsta/cli"
)

type CreateUser struct {
	ui       terminal.UI
	config   configuration.Reader
	userRepo api.UserRepository
}

func NewCreateUser(ui terminal.UI, config configuration.Reader, userRepo api.UserRepository) (cmd CreateUser) {
	cmd.ui = ui
	cmd.config = config
	cmd.userRepo = userRepo
	return
}

func (command CreateUser) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "create-user",
		Description: "Create a new user",
		Usage:       "CF_NAME create-user USERNAME PASSWORD",
	}
}

func (cmd CreateUser) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 2 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c)
	}

	reqs = append(reqs, requirementsFactory.NewLoginRequirement())

	return
}

func (cmd CreateUser) Run(c *cli.Context) {
	username := c.Args()[0]
	password := c.Args()[1]

	cmd.ui.Say("Creating user %s as %s...",
		terminal.EntityNameColor(username),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	err := cmd.userRepo.Create(username, password)
	switch err.(type) {
	case nil:
	case *errors.ModelAlreadyExistsError:
		cmd.ui.Warn("%s", err.Error())
	default:
		cmd.ui.Failed("Error creating user %s.\n%s", terminal.EntityNameColor(username), err.Error())
	}

	cmd.ui.Ok()
	cmd.ui.Say("\nTIP: Assign roles with '%s set-org-role' and '%s set-space-role'", cf.Name(), cf.Name())
}
