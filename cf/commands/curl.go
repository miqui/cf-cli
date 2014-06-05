package commands

import (
	"errors"
	"github.com/starkandwayne/cf-cli/cf/api"
	"github.com/starkandwayne/cf-cli/cf/command_metadata"
	"github.com/starkandwayne/cf-cli/cf/configuration"
	"github.com/starkandwayne/cf-cli/cf/flag_helpers"
	"github.com/starkandwayne/cf-cli/cf/requirements"
	"github.com/starkandwayne/cf-cli/cf/terminal"
	"github.com/starkandwayne/cf-cli/cf/trace"
	"github.com/codegangsta/cli"
	"strings"
)

type Curl struct {
	ui       terminal.UI
	config   configuration.Reader
	curlRepo api.CurlRepository
}

func NewCurl(ui terminal.UI, config configuration.Reader, curlRepo api.CurlRepository) (cmd *Curl) {
	cmd = new(Curl)
	cmd.ui = ui
	cmd.config = config
	cmd.curlRepo = curlRepo
	return
}

func (command *Curl) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "curl",
		Description: "Executes a raw request, content-type set to application/json by default",
		Usage:       "CF_NAME curl PATH [-X METHOD] [-H HEADER] [-d DATA] [-i]",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "X", Value: "GET", Usage: "HTTP method (GET,POST,PUT,DELETE,etc)"},
			flag_helpers.NewStringSliceFlag("H", "Custom headers to include in the request, flag can be specified multiple times"),
			flag_helpers.NewStringFlag("d", "HTTP data to include in the request body"),
			cli.BoolFlag{Name: "i", Usage: "Include response headers in the output"},
			cli.BoolFlag{Name: "v", Usage: "Enable CF_TRACE output for all requests and responses"},
		},
	}
}

func (cmd *Curl) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 1 {
		err = errors.New("Incorrect number of arguments")
		cmd.ui.FailWithUsage(c)
		return
	}

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}
	return
}

func (cmd *Curl) Run(c *cli.Context) {
	path := c.Args()[0]
	method := c.String("X")
	headers := c.StringSlice("H")
	body := c.String("d")
	verbose := c.Bool("v")

	reqHeader := strings.Join(headers, "\n")

	if verbose {
		trace.EnableTrace()
	}

	respHeader, respBody, apiErr := cmd.curlRepo.Request(method, path, reqHeader, body)
	if apiErr != nil {
		cmd.ui.Failed("Error creating request:\n%s", apiErr.Error())
		return
	}

	if verbose {
		return
	}

	if c.Bool("i") {
		cmd.ui.Say("%s", respHeader)
	}

	cmd.ui.Say("%s", respBody)
	return
}
