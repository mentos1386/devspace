package remove

import (
	"github.com/devspace-cloud/devspace/pkg/devspace/config/configutil"
	"github.com/devspace-cloud/devspace/pkg/devspace/configure"
	deployUtil "github.com/devspace-cloud/devspace/pkg/devspace/deploy/util"
	"github.com/devspace-cloud/devspace/pkg/devspace/kubectl"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/devspace-cloud/devspace/pkg/util/stdinutil"
	"github.com/spf13/cobra"
)

type deploymentCmd struct {
	RemoveAll bool
}

func newDeploymentCmd() *cobra.Command {
	cmd := &deploymentCmd{}

	deploymentCmd := &cobra.Command{
		Use:   "deployment [deployment-name]",
		Short: "Removes one or all deployments from devspace configuration",
		Long: `
#######################################################
############ devspace remove deployment ###############
#######################################################
Removes one or all deployments from the devspace
configuration (If you want to delete the deployed 
resources, run 'devspace purge -d deployment_name'):

devspace remove deployment devspace-default
devspace remove deployment --all
#######################################################
	`,
		Args: cobra.MaximumNArgs(1),
		Run:  cmd.RunRemoveDeployment,
	}

	deploymentCmd.Flags().BoolVar(&cmd.RemoveAll, "all", false, "Remove all deployments")

	return deploymentCmd
}

// RunRemoveDeployment executes the specified deployment
func (cmd *deploymentCmd) RunRemoveDeployment(cobraCmd *cobra.Command, args []string) {
	// Set config root
	configExists, err := configutil.SetDevSpaceRoot()
	if err != nil {
		log.Fatal(err)
	}
	if !configExists {
		log.Fatal("Couldn't find any devspace configuration. Please run `devspace init`")
	}

	name := ""
	if len(args) > 0 {
		name = args[0]
	}

	shouldPurgeDeployment := *stdinutil.GetFromStdin(&stdinutil.GetFromStdinParams{
		Question:     "Do you want to delete all deployment resources deployed?",
		DefaultValue: "yes",
		Options: []string{
			"yes",
			"no",
		},
	}) == "yes"
	if shouldPurgeDeployment {
		kubectl, err := kubectl.NewClient()
		if err != nil {
			log.Fatalf("Unable to create new kubectl client: %v", err)
		}

		deployments := []string{}
		if cmd.RemoveAll == false {
			deployments = []string{args[0]}
		}

		deployUtil.PurgeDeployments(kubectl, deployments)
	}

	found, err := configure.RemoveDeployment(cmd.RemoveAll, name)
	if err != nil {
		log.Fatal(err)
	}

	if found {
		log.Donef("Successfully removed deployment %s", args[0])
	} else {
		log.Warnf("Couldn't find deployment %s", args[0])
	}
}
