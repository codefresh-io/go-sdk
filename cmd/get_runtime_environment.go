// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getRuntimeEnvironmentCmd represents the getRuntimeEnvironment command
var getRuntimeEnvironmentCmd = &cobra.Command{
	Use:   "runtime-environment",
	Short: "Get a runtime environment",
	Run: func(cmd *cobra.Command, args []string) {
		client := viper.Get("codefresh")
		codefreshClient, _ := client.(codefresh.Codefresh)
		re, _ := codefreshClient.GetRuntimeEnvironment(args[0])
		fmt.Printf(re.Metadata.Name)
	},
}

func init() {
	getCmd.AddCommand(getRuntimeEnvironmentCmd)
}
