package requirements

import (
	"fmt"
	"github.com/starkandwayne/cf-cli/cf"
	"github.com/starkandwayne/cf-cli/cf/configuration"
	"github.com/starkandwayne/cf-cli/cf/terminal"
)

type TargetedSpaceRequirement struct {
	ui     terminal.UI
	config configuration.Reader
}

func NewTargetedSpaceRequirement(ui terminal.UI, config configuration.Reader) TargetedSpaceRequirement {
	return TargetedSpaceRequirement{ui, config}
}

func (req TargetedSpaceRequirement) Execute() (success bool) {
	if !req.config.HasOrganization() {
		message := fmt.Sprintf("No org and space targeted, use '%s' to target an org and space",
			terminal.CommandColor(cf.Name()+" target -o ORG -s SPACE"))
		req.ui.Failed(message)
		return false
	}

	if !req.config.HasSpace() {
		message := fmt.Sprintf("No space targeted, use '%s' to target a space", terminal.CommandColor("cf target -s"))
		req.ui.Failed(message)
		return false
	}

	return true
}
