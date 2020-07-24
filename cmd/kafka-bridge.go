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
	"k8s.io/kubectl/pkg/cmd/get"
	"time"
)

var (
	kafkaBridgeNamespaceFlag string
)

func kafkaBridgeRoute() {
	co := utils.NewCommandOptions()

	//Setup kafka bridge
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/dmc5179/iot-dev/ocs/yamls/kafka/bridge/route.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/dmc5179/iot-dev/ocs/yamls/kafka/bridge/kafka-bridge.yaml")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext(kafkaBridgeNamespaceFlag)

	//Reload config flags after switching context

	log.Println("Provision Kafka Http Bridge using route")
	for _, command := range co.Commands {
		cmd := apply.NewCmdApply("kubectl", co.CurrentFactory, IOStreams)
		//Kubectl signals missing field, set validate to false to ignore this
		cmd.Flags().Set("validate", "false")
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Print(out.String())
		out.Reset()
	}

	//After pods for Kafka bridge is provisioned wait for them to become ready before moving on
	log.Print("Waiting for Kafka Bridge to be ready")
	podStatus := utils.NewpodStatus()
	for podStatus.Running >= 5 {
		cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
		cmd.Flags().Set("output", "yaml")
		cmd.Run(cmd, []string{"pods"})
		podStatus.CountPods(out.Bytes())
		log.Debug(podStatus)
		log.Info("Waiting for Kafka Bridge...")
		out.Reset()
		time.Sleep(5 * time.Second)
	}
	log.Print("Kafka Deployment is ready")

	cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
	cmd.Flags().Set("output", "jsonpath={.spec.host}")
	cmd.Run(cmd, []string{"route", "my-bridge-route"})
	log.Info("To check status of Kafka HTTP bridge run 'curl -v GET " + out.String() + "/healthy'")
	out.Reset()
}

func kafkaBridge() {

	co := utils.NewCommandOptions()

	//Setup kafka bridge
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/dmc5179/iot-dev/ocs/yamls/kafka/bridge/kafka-bridge.yaml")
	//Setup Nginix Ingress **CONVERT TO OPENSHIFT ROUTE AT SOME POINT** to connect to bridge from outside the cluster
	//Get Nginix controller and apply to cluster
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/mandatory.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/provider/cloud-generic.yaml")
	//Setup the K8s ingress resource

	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/dmc5179/iot-dev/ocs/yamls/kafka/bridge/ingress.yaml")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext(kafkaBridgeNamespaceFlag)

	//Reload config flags after switching context

	log.Println("Provision Kafka Http Bridge")
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

	//After pods for Kafka bridge is provisioned wait for them to become ready before moving on
	log.Print("Waiting for Kafka Bridge to be ready")
	podStatus := utils.NewpodStatus()
	for podStatus.Running != 5 {
		cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
		cmd.Flags().Set("output", "yaml")
		cmd.Run(cmd, []string{"pods"})
		podStatus.CountPods(out.Bytes())
		log.Debug(podStatus)
		log.Info("Waiting for Kafka Bridge...")
		out.Reset()
		time.Sleep(5 * time.Second)
	}
	log.Print("Kafka Deployment is ready")

	cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
	cmd.Flags().Set("output", "jsonpath={.spec.host}")
	cmd.Run(cmd, []string{"pods"})
	log.Info("To check status of Kafka HTTP bridge run 'curl -v GET " + out.String() + "/healthy'")
	out.Reset()

}

// bridgeCmd represents the bridge command
var kafkaBridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Setup Kafka bridge to send data over to the Kafka cluster",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		fstatus, _ := cmd.Flags().GetBool("route")
		if fstatus { // if status is true, call addFloat
			log.Println("Kafka Http Bridge called using Route")
			kafkaBridgeRoute()
		} else {
			log.Println("Kafka Http Bridge called using Ingress")
			kafkaBridge()
		}

	},
}

func init() {
	kafkaCmd.AddCommand(kafkaBridgeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bridgeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bridgeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	kafkaBridgeCmd.Flags().StringVarP(&kafkaBridgeNamespaceFlag, "namespace", "n", "kafka", "Option to specify namespace for kafka deletion, defaults to 'kafka'")
	//Default to using Route
	kafkaBridgeCmd.Flags().BoolP("route", "r", true, "Setup kafka bridge using route, defaults to using ingress")
}
