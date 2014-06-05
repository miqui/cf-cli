package service

import (
	"github.com/starkandwayne/cf-cli/cf/api"
	"github.com/starkandwayne/cf-cli/cf/command_metadata"
	"github.com/starkandwayne/cf-cli/cf/configuration"
	"github.com/starkandwayne/cf-cli/cf/errors"
	"github.com/starkandwayne/cf-cli/cf/requirements"
	"github.com/starkandwayne/cf-cli/cf/terminal"
	"github.com/codegangsta/cli"
)

type DeleteService struct {
	ui                 terminal.UI
	config             configuration.Reader
	serviceRepo        api.ServiceRepository
	serviceInstanceReq requirements.ServiceInstanceRequirement
}

func NewDeleteService(ui terminal.UI, config configuration.Reader, serviceRepo api.ServiceRepository) (cmd *DeleteService) {
	cmd = new(DeleteService)
	cmd.ui = ui
	cmd.config = config
	cmd.serviceRepo = serviceRepo
	return
}

func (command *DeleteService) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "delete-service",
		ShortName:   "ds",
		Description: "Delete a service instance",
		Usage:       "CF_NAME delete-service SERVICE_INSTANCE [-f]",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "f", Usage: "Force deletion without confirmation"},
		},
	}
}

func (cmd *DeleteService) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	var serviceName string

	if len(c.Args()) == 1 {
		serviceName = c.Args()[0]
	}

	if serviceName == "" {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c)
		return
	}

	reqs = []requirements.Requirement{requirementsFactory.NewLoginRequirement()}
	return
}

func (cmd *DeleteService) Run(c *cli.Context) {
	serviceName := c.Args()[0]

	if !c.Bool("f") {
		if !cmd.ui.ConfirmDelete("service", serviceName) {
			return
		}
	}

	cmd.ui.Say("Deleting service %s in org %s / space %s as %s...",
		terminal.EntityNameColor(serviceName),
		terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
		terminal.EntityNameColor(cmd.config.SpaceFields().Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	instance, apiErr := cmd.serviceRepo.FindInstanceByName(serviceName)

	switch apiErr.(type) {
	case nil:
	case *errors.ModelNotFoundError:
		cmd.ui.Ok()
		cmd.ui.Warn("Service %s does not exist.", serviceName)
		return
	default:
		cmd.ui.Failed(apiErr.Error())
		return
	}

	apiErr = cmd.serviceRepo.DeleteService(instance)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
