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

package main

import (
	"fmt"
	"github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/codefresh-io/go-sdk/pkg/utils"
	"os"
)

func main() {
	//cmd.Execute()

	path := fmt.Sprintf("%s/.cfconfig", os.Getenv("HOME"))
	options, err := utils.ReadAuthContext(path, "")
	if err != nil {
		fmt.Println("Failed to read codefresh config file")
		panic(err)
	}
	clientOptions := codefresh.ClientOptions{Host: options.URL,
		Auth: codefresh.AuthOptions{Token: options.Token}}
	cf := codefresh.New(&clientOptions)
	err, ctxts := cf.Contexts().GetGitContexts()
	for _, ctx := range *ctxts {
		fmt.Println(ctx.Metadata.Name)
		_, ctxd := cf.Contexts().GetGitContextByName(ctx.Metadata.Name)
		if ctxd.Spec.Data.Auth.SshPrivateKey != "" {
			fmt.Println(ctxd.Metadata.Name)
		}

	}

}
