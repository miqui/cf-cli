package requirements

import (
	"github.com/starkandwayne/cf-cli/cf/api"
	"github.com/starkandwayne/cf-cli/cf/configuration"
	"github.com/starkandwayne/cf-cli/cf/models"
	"github.com/starkandwayne/cf-cli/cf/terminal"
)

type DomainRequirement interface {
	Requirement
	GetDomain() models.DomainFields
}

type domainApiRequirement struct {
	name       string
	ui         terminal.UI
	config     configuration.Reader
	domainRepo api.DomainRepository
	domain     models.DomainFields
}

func NewDomainRequirement(name string, ui terminal.UI, config configuration.Reader, domainRepo api.DomainRepository) (req *domainApiRequirement) {
	req = new(domainApiRequirement)
	req.name = name
	req.ui = ui
	req.config = config
	req.domainRepo = domainRepo
	return
}

func (req *domainApiRequirement) Execute() bool {
	var apiErr error
	req.domain, apiErr = req.domainRepo.FindByNameInOrg(req.name, req.config.OrganizationFields().Guid)

	if apiErr != nil {
		req.ui.Failed(apiErr.Error())
		return false
	}

	return true
}

func (req *domainApiRequirement) GetDomain() models.DomainFields {
	return req.domain
}
