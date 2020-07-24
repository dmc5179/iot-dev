/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/IoTCLI/cmd/utils"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/delete"
)

var (
	enmasseDestroyNamespaceFlag string
)

func enmasseDestroy() {

	//Make command options for Kafka Setup
	co := utils.NewCommandOptions()

	//Fill ain the commands that must be applied to
	//Install Enmasse Core
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/dmc5179/iot-dev/ocs/yamls/enmasse-infra-namespace.yaml")
	//
	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext(enmasseDestroyNamespaceFlag)

	//Reload config flags after switching context
	log.Println("Destroy Enmasse Messaging Service")
	for _, command := range co.Commands {
		cmd := delete.NewCmdDelete(co.CurrentFactory, IOStreams)
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Print(out.String())
		out.Reset()
	}

}

// destroyCmd represents the destroy command
var enmasseDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy enmasse project",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("destroy called")
		enmasseDestroy()
	},
}

func init() {
	enmasseCmd.AddCommand(enmasseDestroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	enmasseDestroyCmd.Flags().StringVarP(&enmasseDestroyNamespaceFlag, "namespace", "n", "enmasse-infra", "Option to specify namespace for enmasse deployment, defaults to 'enmasse-infra'")
}
