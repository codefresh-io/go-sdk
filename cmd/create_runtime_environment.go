// Copyright © 2018 Codefresh.Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh"

	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setDefault bool

// createRuntimeEnvironmentCmd represents the createRuntimeEnvironment command
var createRuntimeEnvironmentCmd = &cobra.Command{
	Use:     "runtime-environment",
	Aliases: []string{"re"},
	Short:   "Create runtime environment",

	Run: func(cmd *cobra.Command, args []string) {
		client := viper.Get("codefresh")
		codefreshClient := utils.CastToCodefreshOrDie(client)
		cluster := cmd.Flag("cluster").Value.String()
		namespace := cmd.Flag("namespace").Value.String()
		opt := &codefresh.CreateRuntimeOptions{
			Cluster:   cluster,
			Namespace: namespace,
			HasAgent:  true,
		}
		re, err := codefreshClient.RuntimeEnvironments().Create(opt)
		if err == nil {
			fmt.Printf("Runtime-Environment %s created\n", re.Metadata.Name)
		}

		if setDefault {
			fmt.Printf("Setting runtime as default")
			_, err := codefreshClient.RuntimeEnvironments().Default(re.Metadata.Name)
			if err != nil {
				fmt.Printf("Error during setting runtime to be default: %s", err.Error())
			}
			fmt.Printf("Done")
		}
	},
}

func init() {
	createCmd.AddCommand(createRuntimeEnvironmentCmd)
	createRuntimeEnvironmentCmd.Flags().String("cluster", "", "Set name of the cluster (required)")
	createRuntimeEnvironmentCmd.MarkFlagRequired("cluster")
	createRuntimeEnvironmentCmd.Flags().String("namespace", "", "Set name of the namespace (required)")
	createRuntimeEnvironmentCmd.MarkFlagRequired("namespace")
	createRuntimeEnvironmentCmd.Flags().Bool("has-agent", false, "Set if the runtime environment is managed by Codefresh agent")
	createRuntimeEnvironmentCmd.Flags().BoolVar(&setDefault, "set-default", false, "Set the runtime as deault after creation")
}
