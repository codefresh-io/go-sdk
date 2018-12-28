// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getTokensCmd represents the getTokens command
var getTokensCmd = &cobra.Command{
	Use:     "tokens",
	Aliases: []string{"token"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := viper.Get("codefresh")
		codefreshClient, ok := client.(codefresh.Codefresh)
		if !ok {
			panic("Faild to create Codefresh cleint")
		}
		table := internal.CreateTable()
		table.SetHeader([]string{"Created", "ID", "Name", "Reference Subject", "Reference Type", "Token"})
		table.Append([]string{"", "", ""})
		tokens := codefreshClient.GetTokens()
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getTokensCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getTokensCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
