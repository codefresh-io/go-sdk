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
	"time"

	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// waitCmd represents the logs command
var waitCmd = &cobra.Command{
	Use: "wait",
	Run: func(cmd *cobra.Command, args []string) {
		client := viper.Get("codefresh")
		codefreshClient := utils.CastToCodefreshOrDie(client)
		codefreshClient.Workflows().WaitForStatus(args[0], "success", 2*time.Second, 5*time.Minute)
	},
}

func init() {
	rootCmd.AddCommand(waitCmd)
}
