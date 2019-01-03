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

	"github.com/codefresh-io/go-sdk/internal"
	"github.com/codefresh-io/go-sdk/pkg/codefresh"
	humanize "github.com/dustin/go-humanize"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// getPipelineCmd represents the getPipeline command
var getPipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Get pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		client := viper.Get("codefresh")
		codefreshClient, ok := client.(codefresh.Codefresh)
		if !ok {
			internal.DieOnError(fmt.Errorf("Faild to create Codefresh client"))
		}
		table := internal.CreateTable()
		table.SetHeader([]string{"Pipeline Name", "Created At", "Updated At"})
		table.Append([]string{"", "", ""})
		pipelines, err := codefreshClient.GetPipelines()
		internal.DieOnError(err)
		for _, p := range pipelines {
			table.Append([]string{
				p.Metadata.Name,
				humanize.Time(p.Metadata.CreatedAt),
				humanize.Time(p.Metadata.UpdatedAt),
			})
		}
		table.Render()
	},
}

func init() {
	getCmd.AddCommand(getPipelineCmd)
}
