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
	"github.com/codefresh-io/go-sdk/pkg/utils"
	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getTokensCmd = &cobra.Command{
	Use:     "tokens",
	Aliases: []string{"token"},
	Short:   "Get tokens",
	Run: func(cmd *cobra.Command, args []string) {
		client := viper.Get("codefresh")
		codefreshClient := utils.CastToCodefreshOrDie(client)
		table := internal.CreateTable()
		table.SetHeader([]string{"Created", "ID", "Name", "Reference Subject", "Reference Type", "Token"})
		table.Append([]string{"", "", ""})
		tokens, err := codefreshClient.Tokens().Get()
		internal.DieOnError(err)
		for _, t := range tokens {
			table.Append([]string{
				humanize.Time(t.Created),
				t.ID,
				t.Name,
				t.Subject.Ref,
				t.Subject.Type,
				fmt.Sprintf("%s********", t.TokenPrefix),
			})
		}
		table.Render()
	},
}

func init() {
	getCmd.AddCommand(getTokensCmd)
}
