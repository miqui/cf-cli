/*
                       WARNING WARNING WARNING

                Attention all potential contributors

   This testfile is not in the best state. We've been slowly transitioning
   from the built in "testing" package to using Ginkgo. As you can see, we've
   changed the format, but a lot of the setup, test body, descriptions, etc
   are either hardcoded, completely lacking, or misleading.

   For example:

   Describe("Testing with ginkgo"...)      // This is not a great description
   It("TestDoesSoemthing"...)              // This is a horrible description

   Describe("create-user command"...       // Describe the actual object under test
   It("creates a user when provided ..."   // this is more descriptive

   For good examples of writing Ginkgo tests for the cli, refer to

   src/github.com/starkandwayne/cf-cli/cf/commands/application/delete_app_test.go
   src/github.com/starkandwayne/cf-cli/cf/terminal/ui_test.go
   src/github.com/cloudfoundry/loggregator_consumer/consumer_test.go
*/

package domain_test

import (
	. "github.com/starkandwayne/cf-cli/cf/commands/domain"
	"github.com/starkandwayne/cf-cli/cf/configuration"
	testapi "github.com/starkandwayne/cf-cli/testhelpers/api"
	testcmd "github.com/starkandwayne/cf-cli/testhelpers/commands"
	testconfig "github.com/starkandwayne/cf-cli/testhelpers/configuration"
	testreq "github.com/starkandwayne/cf-cli/testhelpers/requirements"
	testterm "github.com/starkandwayne/cf-cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/starkandwayne/cf-cli/testhelpers/matchers"
)

var _ = Describe("Testing with ginkgo", func() {
	It("TestShareDomainRequirements", func() {
		domainRepo := &testapi.FakeDomainRepository{}

		requirementsFactory := &testreq.FakeReqFactory{LoginSuccess: true}
		callShareDomain([]string{"example.com"}, requirementsFactory, domainRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeTrue())

		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: false}
		callShareDomain([]string{"example.com"}, requirementsFactory, domainRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
	})
	It("TestShareDomainFailsWithUsage", func() {

		requirementsFactory := &testreq.FakeReqFactory{LoginSuccess: true}
		domainRepo := &testapi.FakeDomainRepository{}
		ui := callShareDomain([]string{}, requirementsFactory, domainRepo)
		Expect(ui.FailedWithUsage).To(BeTrue())

		ui = callShareDomain([]string{"example.com"}, requirementsFactory, domainRepo)
		Expect(ui.FailedWithUsage).To(BeFalse())
	})
	It("TestShareDomain", func() {

		requirementsFactory := &testreq.FakeReqFactory{LoginSuccess: true}
		domainRepo := &testapi.FakeDomainRepository{}
		ui := callShareDomain([]string{"example.com"}, requirementsFactory, domainRepo)

		Expect(domainRepo.CreateSharedDomainName).To(Equal("example.com"))
		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Creating shared domain", "example.com", "my-user"},
			[]string{"OK"},
		))
	})
})

func callShareDomain(args []string, requirementsFactory *testreq.FakeReqFactory, domainRepo *testapi.FakeDomainRepository) (fakeUI *testterm.FakeUI) {
	fakeUI = new(testterm.FakeUI)
	ctxt := testcmd.NewContext("create-shared-domain", args)
	configRepo := testconfig.NewRepositoryWithAccessToken(configuration.TokenInfo{Username: "my-user"})
	cmd := NewCreateSharedDomain(fakeUI, configRepo, domainRepo)
	testcmd.RunCommand(cmd, ctxt, requirementsFactory)
	return
}
