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
	"github.com/IoTCLI/cmd/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
)

func tensorflowServingSetup() {

	co := utils.NewCommandOptions()

	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/dmc5179/iot-dev/ocs/yamls/tensorflow/setup/tensorflowServing.yaml")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext(tensorflowServingNamespaceFlag)

	//Reload config flags after switching context

	log.Println("Provision Tensorflow Serving Pod")
	for _, command := range co.Commands {
		cmd := apply.NewCmdApply("kubectl", co.CurrentFactory, IOStreams)
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Print(out.String())
		out.Reset()
	}
}

// tensorflowServingSetupCmd represents the tensorflowServingSetup command
var tensorflowServingSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup tensorflow service for analtyics",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("tensorflowServing Setup called")
		tensorflowServingSetup()
	},
}

func init() {
	tensorflowServingCmd.AddCommand(tensorflowServingSetupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tensorflowServingSetupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
