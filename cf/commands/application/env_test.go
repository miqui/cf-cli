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

package application_test

import (
	. "github.com/starkandwayne/cf-cli/cf/commands/application"
	"github.com/starkandwayne/cf-cli/cf/models"
	testcmd "github.com/starkandwayne/cf-cli/testhelpers/commands"
	testconfig "github.com/starkandwayne/cf-cli/testhelpers/configuration"
	testreq "github.com/starkandwayne/cf-cli/testhelpers/requirements"
	testterm "github.com/starkandwayne/cf-cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/starkandwayne/cf-cli/testhelpers/matchers"
)

var _ = Describe("Testing with ginkgo", func() {
	It("TestEnvRequirements", func() {
		requirementsFactory := getEnvDependencies()

		requirementsFactory.LoginSuccess = true
		callEnv([]string{"my-app"}, requirementsFactory)
		Expect(testcmd.CommandDidPassRequirements).To(BeTrue())
		Expect(requirementsFactory.ApplicationName).To(Equal("my-app"))

		requirementsFactory.LoginSuccess = false
		callEnv([]string{"my-app"}, requirementsFactory)
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
	})
	It("TestEnvFailsWithUsage", func() {

		requirementsFactory := getEnvDependencies()
		ui := callEnv([]string{}, requirementsFactory)

		Expect(ui.FailedWithUsage).To(BeTrue())
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
	})
	It("TestEnvListsEnvironmentVariables", func() {

		requirementsFactory := getEnvDependencies()
		requirementsFactory.Application.EnvironmentVars = map[string]string{
			"my-key":  "my-value",
			"my-key2": "my-value2",
		}

		ui := callEnv([]string{"my-app"}, requirementsFactory)

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Getting env variables for app", "my-app", "my-org", "my-space", "my-user"},
			[]string{"OK"},
			[]string{"my-key", "my-value", "my-key2", "my-value2"},
		))
	})

	It("TestEnvShowsEmptyMessage", func() {
		requirementsFactory := getEnvDependencies()
		requirementsFactory.Application.EnvironmentVars = map[string]string{}

		ui := callEnv([]string{"my-app"}, requirementsFactory)

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Getting env variables for app", "my-app"},
			[]string{"OK"},
			[]string{"No env variables exist"},
		))
	})
})

func callEnv(args []string, requirementsFactory *testreq.FakeReqFactory) (ui *testterm.FakeUI) {
	ui = &testterm.FakeUI{}
	ctxt := testcmd.NewContext("env", args)

	configRepo := testconfig.NewRepositoryWithDefaults()
	cmd := NewEnv(ui, configRepo)
	testcmd.RunCommand(cmd, ctxt, requirementsFactory)

	return
}

func getEnvDependencies() (requirementsFactory *testreq.FakeReqFactory) {
	app := models.Application{}
	app.Name = "my-app"
	requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true, Application: app}
	return
}
