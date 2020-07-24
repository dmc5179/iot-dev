/*
Copyright © 2020 RedHat IoT

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
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	//in package import
	"github.com/IoTCLI/cmd/utils"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/get"
	"time"
)

var (
	kafkaSetupNamespaceFlag string
)

func kafkaSetup() {
	//Make command options for Kafka Setup
	co := utils.NewCommandOptions()
	//This section is mimicking the instructions to setup the Strimzi Operator, I.E download the install yaml, and set namespace using sed
	//functionality

	//Make A Temporary file to store output from Sed
	tmpFile, err := ioutil.TempFile(os.TempDir(), "strim-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	defer os.Remove(tmpFile.Name())

	sedCommands := []string{`s/namespace: .*/namespace: kafka/`}

	myOutput := utils.RemoteSed(sedCommands, "https://github.com/strimzi/strimzi-kafka-operator/releases/download/0.17.0/strimzi-cluster-operator-0.17.0.yaml")

	tmpFile.Write(myOutput.Bytes())

	//Close Tempfile after writing
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	//Fill in the commands that must be applied to
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/dmc5179/iot-dev/ocs/yamls/kafka/setup/kafka-namespace.yaml")
	co.Commands = append(co.Commands, tmpFile.Name())
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/dmc5179/iot-dev/ocs/yamls/kafka/setup/kafka.yaml")
	//
	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext(kafkaSetupNamespaceFlag)

	log.Info("Provision Kafka")
	for _, command := range co.Commands {
		cmd := apply.NewCmdApply("kubectl", co.CurrentFactory, IOStreams)
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Info(out.String())
		out.Reset()
	}
	//Remove tempfile when done

	os.Remove(tmpFile.Name())

	//Wait until all resources are up

	//After the pods in openshift-storage are provisioned wait for them to become ready before moving on
	log.Print("Waiting for Kafka deployment to be ready")
	podStatus := utils.NewpodStatus()
	for podStatus.Running != 4 {
		cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
		cmd.Flags().Set("output", "yaml")
		cmd.Run(cmd, []string{"pods"})
		podStatus.CountPods(out.Bytes())
		log.Debug(podStatus)
		log.Info("Waiting for Kafka deployment...")
		out.Reset()
		time.Sleep(5 * time.Second)
	}
	log.Print("Kafka Deployment is ready")

}

// setupCmd represents the setup command
var kafkaSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Kafka with Strimzi Operator on a single Openshift namespace",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Kafka setup called")
		kafkaSetup()

	},
}

func init() {
	kafkaCmd.AddCommand(kafkaSetupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	kafkaSetupCmd.Flags().StringVarP(&kafkaSetupNamespaceFlag, "namespace", "n", "kafka", "Option to specify namespace for kafka deployment, defaults to 'kafka'")
}
