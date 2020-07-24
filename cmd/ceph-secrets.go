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
	"k8s.io/kubectl/pkg/cmd/get"
	"strconv"
)

func getCredentials(user string) {
	co := utils.NewCommandOptions()

	co.Commands = append(co.Commands, "secrets")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	//Switch Context and Reload Config Flags
	co.SwitchContext("openshift-storage")

	log.Print("Get S3 secrets, save for possible later use:")
	cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
	cmd.Flags().Set("output", "jsonpath={.data}")
	cmd.Run(cmd, []string{co.Commands[0], "openshift-storage-object-user-my-store-" + user})
	log.Info(out.String())
	out.Reset()

	log.Print("Ceph Endpoint URL: ")
	cmd = get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
	cmd.Flags().Set("output", "jsonpath={.spec.host}")
	cmd.Run(cmd, []string{"route", "ceph-route"})
	log.Info(out.String())
	out.Reset()
}

// secretsCmd represents the secrets command
var cephSecretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Get S3 secrets from ceph object storage",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Ceph Secrets called")
		cobra.ExactArgs(1)
		if len(args) != 1 {
			log.Fatal("Wrong number of Input arguments expected 1 ceph user got " + strconv.Itoa(len(args)))
		} else {
			getCredentials(args[0])
		}
	},
}

func init() {
	cephCmd.AddCommand(cephSecretsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
