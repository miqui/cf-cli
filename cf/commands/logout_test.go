package commands_test

import (
	"github.com/starkandwayne/cf-cli/cf/commands"
	"github.com/starkandwayne/cf-cli/cf/configuration"
	"github.com/starkandwayne/cf-cli/cf/models"
	testconfig "github.com/starkandwayne/cf-cli/testhelpers/configuration"
	testterm "github.com/starkandwayne/cf-cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("logout command", func() {
	var config configuration.Repository
	BeforeEach(func() {
		org := models.OrganizationFields{}
		org.Name = "MyOrg"

		space := models.SpaceFields{}
		space.Name = "MySpace"

		config = testconfig.NewRepository()
		config.SetAccessToken("MyAccessToken")
		config.SetOrganizationFields(org)
		config.SetSpaceFields(space)
		ui := new(testterm.FakeUI)

		l := commands.NewLogout(ui, config)
		l.Run(nil)
	})

	It("clears access token from the config", func() {
		Expect(config.AccessToken()).To(Equal(""))
	})

	It("clears organization fields from the config", func() {
		Expect(config.OrganizationFields()).To(Equal(models.OrganizationFields{}))
	})

	It("clears space fields from the config", func() {
		Expect(config.SpaceFields()).To(Equal(models.SpaceFields{}))
	})
})
