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
	"os"
	"github.com/spf13/cobra"
)

var (
	customerName  string
	clusterType   string
	outputDir     string
	clusterSize   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mf-templating",
	Short: "Tool for genrating .",
	Long: `Tool for creating a managed kubernetes cluster in either Google Cloud (GKE), AWS (EKS) or Azure (AKS).`,
	Run: func(cmd *cobra.Command, args []string) { },
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
	rootCmd.PersistentFlags().StringVarP(&customerName, "customerName", "c", "", "The name of the customer whom the template is being generated for (required)")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "outputDir", "o", "", "The directory location to output the template(s). (default '../<customerName>/<clusterType>')")
	rootCmd.PersistentFlags().StringVarP(&clusterSize, "clusterSize", "s", "", "Provide the clusterSize(for prod) eg: large, small (for non-prod) enter any value it defaults to c2-standard-8 (required)" )
	rootCmd.MarkPersistentFlagRequired("customerName")
	rootCmd.MarkPersistentFlagRequired("clusterSize")
}