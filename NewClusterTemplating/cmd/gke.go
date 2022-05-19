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

type GKEValues struct {
	CloudProvider string `json:"cloud_provider"`
	ProjectID     string `json:"project_id"`
	ClusterName   string `json:"cluster_name"`
	CustomerName  string `json:"customer_name"`
	ClusterType   string `json:"cluster_type"`
	Region        string `json:"region"`
	Domain        string `json:"domain"`
	MachineType   string `json:"machine_type"`
	OutputDir     string `json:"output_dir"`
	ClusterSize   string `json:"clusterSize"`
}

var (
	projectId string
	region    string
)

var gkeCmd = &cobra.Command{
	Use:   "gke",
	Short: "Generates a terraform script for creating a new k8s cluster in GKE.",
	Long:  "Generates a terraform script for creating a new k8s cluster in GKE.",
	Run: func(cmd *cobra.Command, args []string) {
		gkeExecutor()
	},
}

func init() {
	clusterCmd.AddCommand(gkeCmd)
	gkeCmd.Flags().StringVarP(&clusterType, "clusterType", "t", "", "The type of cluster to generate the templates for. eg: prod, non-prod, prod-<xyz>, non-prod-<xyz> (required)")
	gkeCmd.Flags().StringVarP(&projectId, "projectId", "p", "", "The Google Cloud project Id to create the cluster in (required)")
	gkeCmd.Flags().StringVarP(&region, "region", "r", "", "The region to deploy the cluster to. (required)")
	gkeCmd.Flags().StringVarP(&clusterSize, "clusterSize", "s", "", "The siZe of the cluster large or small for prod(required)")
	gkeCmd.MarkFlagRequired("clusterType")
	gkeCmd.MarkFlagRequired("projectId")
	gkeCmd.MarkFlagRequired("region")
	gkeCmd.MarkFlagRequired("clusterSize")
}

func gkeExecutor() {
	gkeValues := new(GKEValues)
	gkeValues.ProjectID = projectId
	gkeValues.CustomerName = customerName
	gkeValues.ClusterType = clusterType
	gkeValues.ClusterName = fmt.Sprintf("lw-%s-%s", strings.ToLower(customerName), strings.ToLower(clusterType))
	gkeValues.Region = region
	gkeValues.CloudProvider = "gcp"
	gkeValues.ClusterSize = clusterSize

	//if projectId == "managed-fusion" {
	gkeValues.Domain = "b.lucidworks.cloud"
	//} else if projectId == "managed-fusion-dev" {
  //		gkeValues.Domain = "lucidworksdeployments.com"
	//}
	if clusterSize == "small" && clusterType == "prod" {
		gkeValues.MachineType = "c2-standard-4"
	} else if clusterSize == "small" && clusterType == "non-prod" {
		gkeValues.MachineType = "c2-standard-4"
	} else {
		gkeValues.MachineType = "c2-standard-8"
	}

	if outputDir == "" {
		outputDir = fmt.Sprintf("../%s", strings.ToLower(customerName))
	}
	gkeValues.OutputDir = outputDir

	if !dirExists(fmt.Sprintf("%s/repo", outputDir)) {
		templateDirectory(gkeValues, "templates/terraform/repo", fmt.Sprintf("%s/repo", outputDir))
	}
	templateDirectory(gkeValues, "templates/terraform/gke", fmt.Sprintf("%s/%s", outputDir, strings.ToLower(clusterType)))
}
