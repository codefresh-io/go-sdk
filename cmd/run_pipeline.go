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
	"errors"
	"fmt"

	"github.com/codefresh-io/go-sdk/internal"
	"github.com/codefresh-io/go-sdk/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runPipelineCmd represents the runPipeline command
var runPipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"exec"},
	Example: "cfcli run pipeline [pipeline_name_1] [pipeline_name_2] ...",
	Short:   "Run one or more pipelines",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires name of the pipeline")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client := viper.Get("codefresh")
		codefreshClient := utils.CastToCodefreshOrDie(client)
		for _, name := range args {
			build, err := codefreshClient.Pipelines().Run(name)
			internal.DieOnError(err)
			fmt.Printf("Pipeline started with ID: %s\n", build)
		}
	},
}

func init() {
	runCmd.AddCommand(runPipelineCmd)
}
