// Copyright © 2019 Codefresh.Inc
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
	"errors"
	"fmt"

	"github.com/codefresh-io/go-sdk/internal"
	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteRuntimeEnvironmentCmd represents the deleteRuntimeEnvironment command
var deleteRuntimeEnvironmentCmd = &cobra.Command{
	Use:     "runtime-environment",
	Example: "cfcl delete runtime-environment [name_1] [name_2] ...",
	Short:   "Delete a runtime-environment",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires name of the runtime-environment")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client := viper.Get("codefresh")
		codefreshClient := utils.CastToCodefreshOrDie(client)
		for _, name := range args {
			_, err := codefreshClient.RuntimeEnvironments().Delete(name)
			internal.DieOnError(err)
			fmt.Printf("Runtime-environment %s deleted", name)
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteRuntimeEnvironmentCmd)
}
