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
	"os"

	"github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/codefresh-io/go-sdk/pkg/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cfctl",
	Short: "A command line application for Codefresh",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// },

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		configPath := cmd.Flag("cfconfig").Value.String()
		if configPath == "" {
			configPath = fmt.Sprintf("%s/.cfconfig", os.Getenv("HOME"))
		} else {
			fmt.Printf("Nope...\n")
			fmt.Println(configPath)
		}
		context, err := utils.ReadAuthContext(configPath, cmd.Flag("context").Value.String())
		if err != nil {
			return err
		}
		client := codefresh.New(&codefresh.ClietOptions{
			Auth: codefresh.AuthOptions{
				Token: context.Token,
			},
			Host: context.URL,
		})
		viper.Set("codefresh", client)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "cfconfig", "", "config file (default is $HOME/.cfconfig)")
	rootCmd.PersistentFlags().String("context", "", "name of the context from --cfconfig (default is current-context)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cfconfig" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cfconfig")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
