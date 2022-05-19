// Copyright © 2021 Lucidworks
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
	"strings"
	"github.com/spf13/cobra"
)

type REPOValues struct {
	CustomerName string `json:"customer_name"`
	OutputDir    string `json:"output_dir"`
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Generates a terraform script for setting up a github repo for the customer.",
	Long:  "Generates a terraform script for setting up a github repo for the customer.",
	Run: func(cmd *cobra.Command, args []string) {
		repoExecutor()
	},
}

func init() {
	rootCmd.AddCommand(repoCmd)
}

func repoExecutor() {
	repoValues := new(REPOValues)
	repoValues.CustomerName = customerName

	if outputDir == "" {
		outputDir = fmt.Sprintf("../%s", strings.ToLower(customerName))
	}
	repoValues.OutputDir = outputDir

	templateDirectory(repoValues, "templates/terraform/repo", fmt.Sprintf("%s/repo", outputDir))
}
